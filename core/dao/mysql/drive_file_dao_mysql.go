package mysql

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/model"
)

type DriveFileDaoMysql struct {
	db *sql.DB
}

func NewDriveFileDaoMysql(db *sql.DB) *DriveFileDaoMysql {
	return &DriveFileDaoMysql{db: db}
}

func (dao *DriveFileDaoMysql) readById(id string, errorIfNotExists bool) (*model.DriveFile, error) {
	var driveFile model.DriveFile
	row := dao.db.QueryRow("SELECT * FROM backup_drive_files WHERE id = ?", id)
	var createdAt lib.SqlNonNullableTime
	var updatedAt lib.SqlNullableTime
	err := row.Scan(&driveFile.Id, &driveFile.BackupId, &driveFile.Provider, &driveFile.Label, &driveFile.Path, &driveFile.Status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			if errorIfNotExists {
				return nil, fmt.Errorf("no drive files found with id %s", id)
			} else {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to read drive files by id: %v", err)
	}
	driveFile.CreatedAt = createdAt.Time
	if updatedAt.Valid {
		driveFile.UpdatedAt = updatedAt.Time
	}
	return &driveFile, nil
}

func (dao *DriveFileDaoMysql) ReadOrError(id string) (*model.DriveFile, error) {
	return dao.readById(id, true)
}

func (dao *DriveFileDaoMysql) Create(data model.DriveFile) (string, error) {
	id := uuid.NewString()
	_, err := dao.db.Exec("INSERT INTO backup_drive_files (id, backup_id, provider, label, path, status) VALUES (?, ?, ?, ?, ?, ?)", id, data.BackupId, data.Provider, data.Label, data.Path, data.Status)
	if err != nil {
		return "", fmt.Errorf("failed to insert drive file: %v", err)
	}
	return id, nil
}

func (dao *DriveFileDaoMysql) Update(id string, data model.DriveFile) (string, error) {
	_, err := dao.db.Exec("UPDATE backup_drive_files SET backup_id = ?, provider = ?, label = ?, path = ?, status = ? WHERE id = ?", data.BackupId, data.Provider, data.Label, data.Path, data.Status, id)
	if err != nil {
		return "", fmt.Errorf("failed to update drive file: %v", err)
	}
	return id, nil
}

func (dao *DriveFileDaoMysql) Delete(id string) error {
	_, err := dao.db.Exec("DELETE FROM backup_drive_files WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete drive file: %v", err)
	}
	return nil
}
