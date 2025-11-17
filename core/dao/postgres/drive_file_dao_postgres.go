package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/herytz/backupman/core/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DriveFileDaoPostgres struct {
	db *pgxpool.Pool
}

func NewDriveFileDaoPostgres(db *pgxpool.Pool) *DriveFileDaoPostgres {
	return &DriveFileDaoPostgres{db: db}
}

func (dao *DriveFileDaoPostgres) readById(id string, errorIfNotExists bool) (*model.DriveFile, error) {
	var driveFile model.DriveFile
	var updatedAt *time.Time

	err := dao.db.QueryRow(context.Background(), "SELECT id, backup_id, provider, label, path, status, created_at, updated_at FROM backup_drive_files WHERE id = $1", id).Scan(&driveFile.Id, &driveFile.BackupId, &driveFile.Provider, &driveFile.Label, &driveFile.Path, &driveFile.Status, &driveFile.CreatedAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			if errorIfNotExists {
				return nil, fmt.Errorf("no drive files found with id %s", id)
			} else {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to read drive files by id: %v", err)
	}
	if updatedAt != nil {
		driveFile.UpdatedAt = *updatedAt
	}
	return &driveFile, nil
}

func (dao *DriveFileDaoPostgres) ReadOrError(id string) (*model.DriveFile, error) {
	return dao.readById(id, true)
}

func (dao *DriveFileDaoPostgres) Create(data model.DriveFile) (string, error) {
	id := uuid.NewString()
	_, err := dao.db.Exec(context.Background(), "INSERT INTO backup_drive_files (id, backup_id, provider, label, path, status) VALUES ($1, $2, $3, $4, $5, $6)", id, data.BackupId, data.Provider, data.Label, data.Path, data.Status)
	if err != nil {
		return "", fmt.Errorf("failed to insert drive file: %v", err)
	}
	return id, nil
}

func (dao *DriveFileDaoPostgres) Update(id string, data model.DriveFile) (string, error) {
	_, err := dao.db.Exec(context.Background(), "UPDATE backup_drive_files SET backup_id = $1, provider = $2, label = $3, path = $4, status = $5 WHERE id = $6", data.BackupId, data.Provider, data.Label, data.Path, data.Status, id)
	if err != nil {
		return "", fmt.Errorf("failed to update drive file: %v", err)
	}
	return id, nil
}

func (dao *DriveFileDaoPostgres) Delete(id string) error {
	_, err := dao.db.Exec(context.Background(), "DELETE FROM backup_drive_files WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete drive file: %v", err)
	}
	return nil
}
