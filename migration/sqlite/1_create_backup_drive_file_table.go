package sqlite

import (
	"database/sql"
	"fmt"
)

func RunCreateBackupDriveFileTable(cnx *sql.DB) error {
	backupTableQuery := `
CREATE TABLE backups (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    label TEXT NOT NULL,
    dump_path TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
)
	`
	_, err := cnx.Exec(backupTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create backups table => %w", err)
	}

	driveFileTableQuery := `
CREATE TABLE backup_drive_files (
    id TEXT PRIMARY KEY,
    backup_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    label TEXT NOT NULL,
    path TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (backup_id) REFERENCES backups(id) ON DELETE CASCADE
)
	`
	_, err = cnx.Exec(driveFileTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create backup_drive_files table => %w", err)
	}

	return nil
}
