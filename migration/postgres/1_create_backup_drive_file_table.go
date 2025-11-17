package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunCreateBackupDriveFileTable(cnx *pgxpool.Pool) error {
	backupTableQuery := `
CREATE TABLE backups (
    id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    label VARCHAR(50) NOT NULL,
    dump_path VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
)
	`
	_, err := cnx.Exec(context.Background(), backupTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create backups table => %w", err)
	}

	driveFileTableQuery := `
CREATE TABLE backup_drive_files (
    id VARCHAR(36) PRIMARY KEY,
    backup_id VARCHAR(36) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    label VARCHAR(50) NOT NULL,
    path VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (backup_id) REFERENCES backups(id) ON DELETE CASCADE
)
	`
	_, err = cnx.Exec(context.Background(), driveFileTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create drive_files table => %w", err)
	}

	return nil
}
