package postgres

import (
	"context"
	"fmt"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/lib"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	version string
	fn      func(cnx *pgxpool.Pool) error
}

func RunPostgres(db application.PostgresDbConfig) error {
	cnx, err := lib.NewPostgresConnection(db.Host, db.Port, db.User, db.Password, db.Database, db.Tls)
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
		err := cnx.QueryRow(context.Background(), "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'backupman_migrations'").Scan(&tableVersionExists)
		if err != nil {
			return fmt.Errorf("failed to check if migrations table exists => %w", err)
		}

		if tableVersionExists == 0 {
			_, err = cnx.Exec(context.Background(), "CREATE TABLE backupman_migrations (version VARCHAR(255) NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
			if err != nil {
				return fmt.Errorf("failed to create backupman_migrations table => %w", err)
			}
		}

		var version string
		err = cnx.QueryRow(context.Background(), "SELECT version FROM backupman_migrations WHERE version = $1", migration.version).Scan(&version)
		if err != nil && err.Error() != "no rows in result set" {
			return fmt.Errorf("failed to check migration version => %w", err)
		}

		if version == migration.version {
			continue
		}

		err = migration.fn(cnx)
		if err != nil {
			return err
		}

		_, err = cnx.Exec(context.Background(), "INSERT INTO backupman_migrations (version) VALUES ($1)", migration.version)
		if err != nil {
			return fmt.Errorf("failed to insert backupman_migrations version => %w", err)
		}
	}

	return nil
}
