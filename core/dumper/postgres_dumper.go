package dumper

import (
	"context"
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
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDumper struct {
	Label     string
	TmpFolder string
	db        *pgxpool.Pool
}

type postgresTable struct {
	Name   string
	SQL    string
	Values string
}

type postgresDump struct {
	DumpVersion   string
	ServerVersion string
	Tables        []*postgresTable
	CompleteTime  string
}

const postgresVersion = "0.1.0"

const postgresTmpl = `-- Backupmap PostgreSQL Dump {{ .DumpVersion }}
--
-- ------------------------------------------------------
-- Server version	{{ .ServerVersion }}

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

{{range .Tables}}
--
-- Table structure for table "{{ .Name }}"
--

DROP TABLE IF EXISTS "{{ .Name }}" CASCADE;
{{ .SQL }};

--
-- Dumping data for table "{{ .Name }}"
--

{{ if .Values }}
{{ .Values }}
{{ end }}
{{ end }}

-- Dump completed on {{ .CompleteTime }}
`

func NewPostgresDumper(label, tmpFolder, host string, port int, user, password, database string, tls bool) *PostgresDumper {
	db, err := lib.NewPostgresConnection(host, port, user, password, database, tls)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	postgresDumper := &PostgresDumper{
		db:        db,
		Label:     label,
		TmpFolder: tmpFolder,
	}
	postgresDumper.setup()
	return postgresDumper
}

func (p *PostgresDumper) Dump() (string, error) {
	filename := uuid.NewString() + ".sql"
	filenamePath := path.Join(p.TmpFolder, filename)

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

	data := postgresDump{
		DumpVersion: postgresVersion,
		Tables:      make([]*postgresTable, 0),
	}

	data.ServerVersion, err = p.getServerVersion()
	if err != nil {
		return "", err
	}

	tables, err := p.getTables()
	if err != nil {
		return "", err
	}

	for _, name := range tables {
		t, err := p.createTable(name)
		if err != nil {
			return "", err
		}
		data.Tables = append(data.Tables, t)
	}

	data.CompleteTime = time.Now().String()

	tm, err := template.New("postgresdump").Parse(postgresTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}
	err = tm.Execute(file, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %s", err)
	}

	return filenamePath, nil
}

func (p *PostgresDumper) getServerVersion() (string, error) {
	var serverVersion string
	err := p.db.QueryRow(context.Background(), "SELECT version()").Scan(&serverVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get server version %s", err)
	}
	return serverVersion, nil
}

func (p *PostgresDumper) getTables() ([]string, error) {
	tables := make([]string, 0)
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`
	rows, err := p.db.Query(context.Background(), query)
	if err != nil {
		return tables, fmt.Errorf("failed to get tables: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return tables, fmt.Errorf("failed to scan row: %s", err)
		}
		tables = append(tables, tableName)
	}
	return tables, rows.Err()
}

func (p *PostgresDumper) createTable(name string) (*postgresTable, error) {
	var err error
	t := &postgresTable{Name: name}

	t.SQL, err = p.createTableSQL(name)
	if err != nil {
		return t, err
	}

	t.Values, err = p.createTableValues(name)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (p *PostgresDumper) createTableSQL(name string) (string, error) {
	// Get column definitions
	columnsQuery := `
		SELECT
			column_name,
			data_type,
			character_maximum_length,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position
	`
	rows, err := p.db.Query(context.Background(), columnsQuery, name)
	if err != nil {
		return "", fmt.Errorf("failed to get columns for table %s: %s", name, err)
	}
	defer rows.Close()

	var columnDefs []string
	for rows.Next() {
		var columnName, dataType, isNullable string
		var charMaxLength *int
		var columnDefault *string

		err := rows.Scan(&columnName, &dataType, &charMaxLength, &isNullable, &columnDefault)
		if err != nil {
			return "", fmt.Errorf("failed to scan column info: %s", err)
		}

		// Build column definition
		colDef := fmt.Sprintf("\"%s\" %s", columnName, p.mapDataType(dataType, charMaxLength))

		if isNullable == "NO" {
			colDef += " NOT NULL"
		}

		if columnDefault != nil {
			colDef += fmt.Sprintf(" DEFAULT %s", *columnDefault)
		}

		columnDefs = append(columnDefs, colDef)
	}

	if err = rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating columns: %s", err)
	}

	// Get primary key constraint
	pkQuery := `
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = $1::regclass AND i.indisprimary
		ORDER BY a.attnum
	`
	pkRows, err := p.db.Query(context.Background(), pkQuery, name)
	if err != nil {
		return "", fmt.Errorf("failed to get primary key for table %s: %s", name, err)
	}
	defer pkRows.Close()

	var pkColumns []string
	for pkRows.Next() {
		var colName string
		if err := pkRows.Scan(&colName); err != nil {
			return "", fmt.Errorf("failed to scan primary key column: %s", err)
		}
		pkColumns = append(pkColumns, fmt.Sprintf("\"%s\"", colName))
	}

	if len(pkColumns) > 0 {
		pkDef := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ", "))
		columnDefs = append(columnDefs, pkDef)
	}

	createSQL := fmt.Sprintf("CREATE TABLE \"%s\" (\n  %s\n)", name, strings.Join(columnDefs, ",\n  "))
	return createSQL, nil
}

func (p *PostgresDumper) mapDataType(dataType string, charMaxLength *int) string {
	switch dataType {
	case "character varying":
		if charMaxLength != nil {
			return fmt.Sprintf("VARCHAR(%d)", *charMaxLength)
		}
		return "VARCHAR"
	case "character":
		if charMaxLength != nil {
			return fmt.Sprintf("CHAR(%d)", *charMaxLength)
		}
		return "CHAR"
	case "timestamp without time zone":
		return "TIMESTAMP"
	case "timestamp with time zone":
		return "TIMESTAMPTZ"
	case "time without time zone":
		return "TIME"
	case "time with time zone":
		return "TIMETZ"
	default:
		return strings.ToUpper(dataType)
	}
}

func (p *PostgresDumper) createTableValues(name string) (string, error) {
	query := fmt.Sprintf("SELECT * FROM \"%s\"", name)
	rows, err := p.db.Query(context.Background(), query)
	if err != nil {
		return "", fmt.Errorf("cannot get table %s values: %s", name, err)
	}
	defer rows.Close()

	// Get column names and count
	fieldDescriptions := rows.FieldDescriptions()
	if len(fieldDescriptions) == 0 {
		return "", nil
	}

	var insertStatements []string
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %s", err)
		}

		dataStrings := make([]string, len(values))
		for i, value := range values {
			if value == nil {
				dataStrings[i] = "NULL"
			} else {
				switch v := value.(type) {
				case string:
					// Escape single quotes by doubling them
					escaped := strings.ReplaceAll(v, "'", "''")
					dataStrings[i] = fmt.Sprintf("'%s'", escaped)
				case []byte:
					// Handle byte arrays as strings
					escaped := strings.ReplaceAll(string(v), "'", "''")
					dataStrings[i] = fmt.Sprintf("'%s'", escaped)
				case bool:
					if v {
						dataStrings[i] = "TRUE"
					} else {
						dataStrings[i] = "FALSE"
					}
				case time.Time:
					dataStrings[i] = fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
				default:
					dataStrings[i] = fmt.Sprintf("%v", v)
				}
			}
		}

		insertStatements = append(insertStatements, fmt.Sprintf("INSERT INTO \"%s\" VALUES (%s);", name, strings.Join(dataStrings, ", ")))
	}

	if err = rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %s", err)
	}

	return strings.Join(insertStatements, "\n"), nil
}

func (p *PostgresDumper) Health() error {
	return lib.NewHealthPostgres(p.db).Check()
}

func (p *PostgresDumper) GetLabel() string {
	return p.Label
}

func (p *PostgresDumper) setup() {
	err := os.MkdirAll(p.TmpFolder, 0755)
	if err != nil {
		log.Fatalf("failed to setup PostgresDumper (%s) tmpFolder (%s). %s", p.Label, p.TmpFolder, err)
	}
}
