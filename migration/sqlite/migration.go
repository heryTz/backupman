package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/lib"
)

type Migration struct {
	version string
	fn      func(sql *sql.DB) error
}

func RunSqlite(db application.SqliteDbConfig) error {
	cnx, err := lib.NewSqliteConnection(db.DbPath)
	if err != nil {
		return err
	}
	defer cnx.Close()

	migrations := []Migration{
		{
			version: "1",
			fn:      RunCreateBackupDriveFileTable,
		},
	}

	for _, migration := range migrations {
		var tableVersionExists int
		err := cnx.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='backupman_migrations'").Scan(&tableVersionExists)
		if err != nil {
			return fmt.Errorf("failed to check if migrations table exists => %w", err)
		}

		if tableVersionExists == 0 {
			_, err = cnx.Exec("CREATE TABLE backupman_migrations (version TEXT NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
			if err != nil {
				return fmt.Errorf("failed to create backupman_migrations table => %w", err)
			}
		}

		var version string
		err = cnx.QueryRow("SELECT version FROM backupman_migrations WHERE version = ?", migration.version).Scan(&version)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to check migration version => %w", err)
		}

		if version == migration.version {
			continue
		}

		err = migration.fn(cnx)
		if err != nil {
			return err
		}

		_, err = cnx.Exec("INSERT INTO backupman_migrations (version) VALUES (?)", migration.version)
		if err != nil {
			return fmt.Errorf("failed to insert backupman_migrations version => %w", err)
		}
	}

	return nil
}
