package dumper

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/herytz/backupman/core/lib"
)

type MysqlDumper struct {
	Label     string
	TmpFolder string
	db        *sql.DB
}

// Take from: https://github.com/JamesStewy/go-mysqldump

type table struct {
	Name   string
	SQL    string
	Values string
}

type dump struct {
	DumpVersion   string
	ServerVersion string
	Tables        []*table
	CompleteTime  string
}

const version = "0.1.0"

const tmpl = `-- Backupmap SQL Dump {{ .DumpVersion }}
--
-- ------------------------------------------------------
-- Server version	{{ .ServerVersion }}

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
SET NAMES utf8mb4;

{{range .Tables}}
--
-- Table structure for table {{ .Name }} 
--

DROP TABLE IF EXISTS {{ .Name }};
{{ .SQL }};

--
-- Dumping data for table {{ .Name }} 
--

{{ if .Values }}
INSERT INTO {{ .Name }} VALUES {{ .Values }};
{{ end }}
{{ end }}

-- Dump completed on {{ .CompleteTime }}
`

func NewMysqlDumper(label, tmpFolder, host string, port int, user, password, database, tls string) *MysqlDumper {
	db, err := lib.NewConnection(host, port, user, password, database, tls)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	mysqlDumper := &MysqlDumper{
		db:        db,
		Label:     label,
		TmpFolder: tmpFolder,
	}
	mysqlDumper.setup()
	return mysqlDumper
}

func (m *MysqlDumper) Dump() (string, error) {
	filename := uuid.NewString() + ".sql"
	filenamePath := path.Join(m.TmpFolder, filename)

	_, err := os.Stat(filenamePath)
	canCreateFile := false
	if err != nil {
		switch err.(type) {
		case *fs.PathError:
			canCreateFile = true
		default:
			return "", fmt.Errorf("failed to check if dump filename (%s) already exists or not: %s", filenamePath, err)
		}
	}

	if !canCreateFile {
		return "", fmt.Errorf("dump filename (%s) already exists", filenamePath)
	}

	file, err := os.Create(filenamePath)
	if err != nil {
		return "", fmt.Errorf("cannot create dump filename (%s): %s", filenamePath, err)
	}
	defer file.Close()

	data := dump{
		DumpVersion: version,
		Tables:      make([]*table, 0),
	}

	data.ServerVersion, err = m.getServerVersion()
	if err != nil {
		return "", err
	}

	tables, err := m.getTables("BASE TABLE")
	if err != nil {
		return "", err
	}

	for _, name := range tables {
		t, err := m.createTable(name, "BASE TABLE")
		if err != nil {
			return "", err
		}
		data.Tables = append(data.Tables, t)
	}

	views, err := m.getTables("VIEW")
	if err != nil {
		return "", err
	}

	for _, name := range views {
		t, err := m.createTable(name, "VIEW")
		if err != nil {
			return "", err
		}
		data.Tables = append(data.Tables, t)
	}

	data.CompleteTime = time.Now().String()

	tm, err := template.New("mysqldump").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}
	err = tm.Execute(file, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %s", err)
	}

	return filenamePath, nil
}

func (m *MysqlDumper) getServerVersion() (string, error) {
	var serverVersion sql.NullString
	err := m.db.QueryRow("SELECT version()").Scan(&serverVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get server version %s", err)
	}
	return serverVersion.String, nil
}

func (m *MysqlDumper) getTables(tableType string) ([]string, error) {
	tables := make([]string, 0)
	rows, err := m.db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = ?", tableType)
	if err != nil {
		return tables, fmt.Errorf("failed to show tables %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table sql.NullString
		err := rows.Scan(&table)
		if err != nil {
			return tables, fmt.Errorf("failed to scan row %s", err)
		}
		tables = append(tables, table.String)
	}
	return tables, rows.Err()
}

func (m *MysqlDumper) createTable(name, tableType string) (*table, error) {
	var err error
	t := &table{Name: name}

	t.SQL, err = m.createTableSQL(name, tableType)
	if err != nil {
		return t, err
	}

	if tableType == "BASE TABLE" {
		t.Values, err = m.createTableValues(name)
		if err != nil {
			return t, err
		}
	}

	return t, nil
}

func (m *MysqlDumper) createTableSQL(name string, tableType string) (string, error) {
	q := "SHOW CREATE TABLE " + name
	if tableType == "VIEW" {
		q = "SHOW CREATE VIEW " + name
	}
	rows, err := m.db.Query(q)
	if err != nil {
		return "", fmt.Errorf("failed to show create table (%s) => %s", name, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("failed to get columns from create table (%s) => %s", name, err)
	}

	if !rows.Next() {
		return "", fmt.Errorf("no rows returned from create table (%s)", name)
	}

	var tableName, tableSql sql.NullString
	var columnScans = make([]sql.NullString, len(columns))
	columnPtrs := make([]interface{}, len(columns))
	columnPtrs[0] = &tableName
	columnPtrs[1] = &tableSql
	for i := 2; i < len(columns); i++ {
		columnPtrs[i] = &columnScans[i]
	}

	err = rows.Scan(columnPtrs...)
	if err != nil {
		return "", fmt.Errorf("failed to scan create table (%s) => %s", name, err)
	}
	if tableName.String != name {
		return "", fmt.Errorf("returned table is not the same as requested table : expected=%s, got=%s", name, tableName.String)
	}
	return tableSql.String, nil
}

func (m *MysqlDumper) createTableValues(name string) (string, error) {
	rows, err := m.db.Query("SELECT * FROM " + name)
	if err != nil {
		return "", fmt.Errorf("cannot get table %s values", name)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("cannot get columns from table %s: %s", name, err)
	}
	if len(columns) == 0 {
		return "", fmt.Errorf("table %s has no columns", name)
	}

	values := make([]string, 0)
	for rows.Next() {
		// data will store the values of each column
		data := make([]*sql.NullString, len(columns))
		// Scan need a pointer to work so we create ptrs to store data pointers
		ptrs := make([]interface{}, len(columns))
		for i := range data {
			ptrs[i] = &data[i]
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			return "", fmt.Errorf("failed to scan row %s", err)
		}

		dataStrings := make([]string, len(columns))
		for i, value := range data {
			if value != nil && value.Valid {
				escaped := strings.ReplaceAll(value.String, "'", "''")
				dataStrings[i] = "'" + escaped + "'"
			} else {
				dataStrings[i] = "NULL"
			}
		}

		values = append(values, "("+strings.Join(dataStrings, ",")+")")
	}

	return strings.Join(values, ","), rows.Err()
}

func (m *MysqlDumper) GetLabel() string {
	return m.Label
}

func (m *MysqlDumper) setup() {
	err := os.MkdirAll(m.TmpFolder, 0755)
	if err != nil {
		log.Fatalf("failed to setup MysqlDumper (%s) tmpFolder (%s). %s", m.Label, m.TmpFolder, err)
	}
}
