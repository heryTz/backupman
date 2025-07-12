package mysql

import (
	"database/sql"
	"fmt"
)

func RunCreateBackupDriveFileTable(cnx *sql.DB) error {
	backupTableQuery := `
CREATE TABLE backups (
    id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    label VARCHAR(50) NOT NULL,
    dump_path VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`
	_, err := cnx.Exec(backupTableQuery)
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
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (backup_id) REFERENCES backups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err = cnx.Exec(driveFileTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create drive_files table => %w", err)
	}

	return nil
}
