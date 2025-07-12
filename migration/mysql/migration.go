package mysql

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

func RunMysql(db application.MysqlDbConfig) error {
	cnx, err := lib.NewConnection(db.Host, db.Port, db.User, db.Password, db.Database, db.Tls)
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
		err := cnx.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'migrations'", db.Database).Scan(&tableVersionExists)
		if err != nil {
			return fmt.Errorf("failed to check if migrations table exists => %w", err)
		}

		if tableVersionExists == 0 {
			_, err = cnx.Exec("CREATE TABLE `migrations` (`version` varchar(255) NOT NULL,`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP)")
			if err != nil {
				return fmt.Errorf("failed to create migrations table => %w", err)
			}
		}

		var version string
		err = cnx.QueryRow("SELECT version FROM migrations WHERE version = ?", migration.version).Scan(&version)
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

		_, err = cnx.Exec("INSERT INTO migrations (version) VALUES (?)", migration.version)
		if err != nil {
			return fmt.Errorf("failed to insert migration version => %w", err)
		}
	}

	return nil
}
