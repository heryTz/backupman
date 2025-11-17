package dumper

import (
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/herytz/backupman/core/lib"
)

type SqliteDumper struct {
	Label     string
	TmpFolder string
	DbPath    string
	db        *sql.DB
}

func NewSqliteDumper(label, tmpFolder, dbPath string) *SqliteDumper {
	db, err := lib.NewSqliteConnection(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	sqliteDumper := &SqliteDumper{
		db:        db,
		Label:     label,
		TmpFolder: tmpFolder,
		DbPath:    dbPath,
	}
	sqliteDumper.setup()
	return sqliteDumper
}

func (s *SqliteDumper) Dump() (string, error) {
	filename := uuid.NewString() + ".db"
	filenamePath := path.Join(s.TmpFolder, filename)

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

	sourceFile, err := os.Open(s.DbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source database file (%s): %s", s.DbPath, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(filenamePath)
	if err != nil {
		return "", fmt.Errorf("cannot create dump filename (%s): %s", filenamePath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy database file: %s", err)
	}

	err = destFile.Sync()
	if err != nil {
		return "", fmt.Errorf("failed to sync file to disk: %s", err)
	}

	return filenamePath, nil
}

func (s *SqliteDumper) Health() error {
	return lib.NewHealthSqlite(s.db).Check()
}

func (s *SqliteDumper) GetLabel() string {
	return s.Label
}

func (s *SqliteDumper) setup() {
	err := os.MkdirAll(s.TmpFolder, 0755)
	if err != nil {
		log.Fatalf("failed to setup SqliteDumper (%s) tmpFolder (%s). %s", s.Label, s.TmpFolder, err)
	}
}
