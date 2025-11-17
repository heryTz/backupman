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

type BackupDaoPostgres struct {
	db *pgxpool.Pool
}

func NewBackupDaoPostgres(db *pgxpool.Pool) *BackupDaoPostgres {
	return &BackupDaoPostgres{db: db}
}

func (dao *BackupDaoPostgres) readById(id string, errorIfNotExists bool) (*model.Backup, error) {
	var backup model.Backup
	var dumpPath *string
	var updatedAt *time.Time

	err := dao.db.QueryRow(context.Background(), "SELECT id, status, label, dump_path, created_at, updated_at FROM backups WHERE id = $1", id).Scan(&backup.Id, &backup.Status, &backup.Label, &dumpPath, &backup.CreatedAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			if errorIfNotExists {
				return nil, fmt.Errorf("no backup found with id %s", id)
			} else {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to read backup by id => %s", err)
	}
	if dumpPath != nil {
		backup.DumpPath = *dumpPath
	}
	if updatedAt != nil {
		backup.UpdatedAt = *updatedAt
	}
	return &backup, nil
}

func (dao *BackupDaoPostgres) ReadOrError(id string) (*model.Backup, error) {
	return dao.readById(id, true)
}

func (dao *BackupDaoPostgres) Create(data model.Backup) (string, error) {
	id := uuid.NewString()
	_, err := dao.db.Exec(context.Background(), "INSERT INTO backups (id, status, label, dump_path) VALUES ($1, $2, $3, $4)", id, data.Status, data.Label, data.DumpPath)
	if err != nil {
		return "", fmt.Errorf("failed to insert backup => %s", err)
	}
	return id, nil
}

func (dao *BackupDaoPostgres) Update(id string, data model.Backup) (string, error) {
	_, err := dao.db.Exec(context.Background(), "UPDATE backups SET status = $1, label = $2, dump_path = $3 WHERE id = $4", data.Status, data.Label, data.DumpPath, id)
	if err != nil {
		return "", fmt.Errorf("failed to update backup => %s", err)
	}
	return id, nil
}

func (dao *BackupDaoPostgres) scanFullBackup(rows pgx.Rows) (map[string]*model.BackupFull, error) {
	results := make(map[string]*model.BackupFull)
	defer rows.Close()

	for rows.Next() {
		var backupScan struct {
			Id        string
			Status    string
			Label     string
			DumpPath  *string
			CreatedAt time.Time
			UpdatedAt *time.Time
		}

		var driveFileScan struct {
			Id        *string
			BackupId  *string
			Provider  *string
			Label     *string
			Path      *string
			Status    *string
			CreatedAt *time.Time
			UpdatedAt *time.Time
		}

		err := rows.Scan(
			&backupScan.Id,
			&backupScan.Status,
			&backupScan.Label,
			&backupScan.DumpPath,
			&backupScan.CreatedAt,
			&backupScan.UpdatedAt,
			&driveFileScan.Id,
			&driveFileScan.BackupId,
			&driveFileScan.Provider,
			&driveFileScan.Label,
			&driveFileScan.Path,
			&driveFileScan.Status,
			&driveFileScan.CreatedAt,
			&driveFileScan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan backup full => %s", err)
		}

		backupFull := model.BackupFull{
			Id:        backupScan.Id,
			Status:    backupScan.Status,
			Label:     backupScan.Label,
			CreatedAt: backupScan.CreatedAt,
		}

		if backupScan.DumpPath != nil {
			backupFull.DumpPath = *backupScan.DumpPath
		}
		if backupScan.UpdatedAt != nil {
			backupFull.UpdatedAt = *backupScan.UpdatedAt
		}

		if results[backupFull.Id] == nil {
			results[backupFull.Id] = &backupFull
		}

		if driveFileScan.Id != nil {
			driveFile := &model.DriveFile{}
			if driveFileScan.Id != nil {
				driveFile.Id = *driveFileScan.Id
			}
			if driveFileScan.BackupId != nil {
				driveFile.BackupId = *driveFileScan.BackupId
			}
			if driveFileScan.Provider != nil {
				driveFile.Provider = *driveFileScan.Provider
			}
			if driveFileScan.Label != nil {
				driveFile.Label = *driveFileScan.Label
			}
			if driveFileScan.Path != nil {
				driveFile.Path = *driveFileScan.Path
			}
			if driveFileScan.Status != nil {
				driveFile.Status = *driveFileScan.Status
			}
			if driveFileScan.CreatedAt != nil {
				driveFile.CreatedAt = *driveFileScan.CreatedAt
			}
			if driveFileScan.UpdatedAt != nil {
				driveFile.UpdatedAt = *driveFileScan.UpdatedAt
			}
			results[backupFull.Id].DriveFiles = append(results[backupFull.Id].DriveFiles, driveFile)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows => %s", err)
	}

	return results, nil
}

func (dao *BackupDaoPostgres) ReadFullById(id string) (*model.BackupFull, error) {
	rows, err := dao.db.Query(context.Background(), "SELECT b.id, b.status, b.label, b.dump_path, b.created_at, b.updated_at, df.id, df.backup_id, df.provider, df.label, df.path, df.status, df.created_at, df.updated_at FROM backups b LEFT JOIN backup_drive_files df ON b.id = df.backup_id WHERE b.id = $1 ORDER BY b.created_at DESC", id)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup full by id %s => %v", id, err)
	}
	results, err := dao.scanFullBackup(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no backup found with id %s", id)
	}
	return results[id], nil
}

func (dao *BackupDaoPostgres) ReadAllFull() ([]model.BackupFull, error) {
	rows, err := dao.db.Query(context.Background(), "SELECT b.id, b.status, b.label, b.dump_path, b.created_at, b.updated_at, df.id, df.backup_id, df.provider, df.label, df.path, df.status, df.created_at, df.updated_at FROM backups b LEFT JOIN backup_drive_files df ON b.id = df.backup_id ORDER BY b.created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to read backup all full => %v", err)
	}
	results, err := dao.scanFullBackup(rows)
	if err != nil {
		return nil, err
	}
	backups := make([]model.BackupFull, 0, len(results))
	for _, backup := range results {
		backups = append(backups, *backup)
	}
	return backups, nil
}

func (dao *BackupDaoPostgres) ReadOlderThan(date time.Time) ([]model.BackupFull, error) {
	rows, err := dao.db.Query(context.Background(), "SELECT b.id, b.status, b.label, b.dump_path, b.created_at, b.updated_at, df.id, df.backup_id, df.provider, df.label, df.path, df.status, df.created_at, df.updated_at FROM backups b LEFT JOIN backup_drive_files df ON b.id = df.backup_id WHERE b.created_at < $1 ORDER BY b.created_at DESC", date)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup older than %s => %v", date, err)
	}
	results, err := dao.scanFullBackup(rows)
	if err != nil {
		return nil, err
	}
	backups := make([]model.BackupFull, 0, len(results))
	for _, backup := range results {
		backups = append(backups, *backup)
	}
	return backups, nil
}

func (dao *BackupDaoPostgres) Delete(id string) error {
	_, err := dao.db.Exec(context.Background(), "DELETE FROM backups WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete backup => %s", err)
	}
	return nil
}
