package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/model"
)

type BackupDaoSqlite struct {
	db *sql.DB
}

func NewBackupDaoSqlite(db *sql.DB) *BackupDaoSqlite {
	return &BackupDaoSqlite{db: db}
}

func (dao *BackupDaoSqlite) readById(id string, errorIfNotExists bool) (*model.Backup, error) {
	var backup model.Backup
	row := dao.db.QueryRow("SELECT * FROM backups WHERE id = ?", id)
	var createdAt lib.SqlNonNullableTime
	var updatedAt lib.SqlNullableTime
	err := row.Scan(&backup.Id, &backup.Status, &backup.Label, &backup.DumpPath, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			if errorIfNotExists {
				return nil, fmt.Errorf("no backup found with id %s", id)
			} else {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to read backup by id => %s", err)
	}
	backup.CreatedAt = createdAt.Time
	if updatedAt.Valid {
		backup.UpdatedAt = updatedAt.Time
	}
	return &backup, nil
}

func (dao *BackupDaoSqlite) ReadOrError(id string) (*model.Backup, error) {
	return dao.readById(id, true)
}

func (dao *BackupDaoSqlite) Create(data model.Backup) (string, error) {
	id := uuid.NewString()
	_, err := dao.db.Exec("INSERT INTO backups (id, status, label, dump_path) VALUES (?, ?, ?, ?)", id, data.Status, data.Label, data.DumpPath)
	if err != nil {
		return "", fmt.Errorf("failed to insert backup => %s", err)
	}
	return id, nil
}

func (dao *BackupDaoSqlite) Update(id string, data model.Backup) (string, error) {
	_, err := dao.db.Exec("UPDATE backups SET status = ?, label = ?, dump_path = ? WHERE id = ?", data.Status, data.Label, data.DumpPath, id)
	if err != nil {
		return "", fmt.Errorf("failed to update backup => %s", err)
	}
	return id, nil
}

func (dao *BackupDaoSqlite) scanFullBackup(rows *sql.Rows) (map[string]*model.BackupFull, error) {
	results := make(map[string]*model.BackupFull)
	for rows.Next() {
		var backupScan struct {
			Id        string
			Status    string
			Label     string
			DumpPath  sql.NullString
			CreatedAt lib.SqlNonNullableTime
			UpdatedAt lib.SqlNullableTime
		}

		var driveFileScan struct {
			Id        sql.NullString
			BackupId  sql.NullString
			Provider  sql.NullString
			Label     sql.NullString
			Path      sql.NullString
			Status    sql.NullString
			CreatedAt lib.SqlNullableTime
			UpdatedAt lib.SqlNullableTime
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
			DumpPath:  backupScan.DumpPath.String,
			CreatedAt: backupScan.CreatedAt.Time,
			UpdatedAt: backupScan.UpdatedAt.Time,
		}

		if results[backupFull.Id] == nil {
			results[backupFull.Id] = &backupFull
		}

		results[backupFull.Id].DriveFiles = append(
			results[backupFull.Id].DriveFiles,
			&model.DriveFile{
				Id:        driveFileScan.Id.String,
				BackupId:  driveFileScan.BackupId.String,
				Provider:  driveFileScan.Provider.String,
				Label:     driveFileScan.Label.String,
				Path:      driveFileScan.Path.String,
				Status:    driveFileScan.Status.String,
				CreatedAt: driveFileScan.CreatedAt.Time,
				UpdatedAt: driveFileScan.UpdatedAt.Time,
			},
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows => %s", err)
	}

	return results, nil
}

func (dao *BackupDaoSqlite) ReadFullById(id string) (*model.BackupFull, error) {
	rows, err := dao.db.Query("SELECT * FROM backups LEFT JOIN backup_drive_files ON backups.id = backup_drive_files.backup_id WHERE backups.id = ? ORDER BY backups.created_at DESC", id)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup full by id %s => %v", id, err)
	}
	defer rows.Close()
	results, err := dao.scanFullBackup(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no backup found with id %s", id)
	}
	return results[id], nil
}

func (dao *BackupDaoSqlite) ReadAllFull() ([]model.BackupFull, error) {
	rows, err := dao.db.Query("SELECT * FROM backups LEFT JOIN backup_drive_files ON backups.id = backup_drive_files.backup_id ORDER BY backups.created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to read backup all full => %v", err)
	}
	defer rows.Close()
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

func (dao *BackupDaoSqlite) ReadOlderThan(date time.Time) ([]model.BackupFull, error) {
	rows, err := dao.db.Query("SELECT * FROM backups LEFT JOIN backup_drive_files ON backups.id = backup_drive_files.backup_id WHERE backups.created_at < ? ORDER BY backups.created_at DESC", date)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup older than %s => %v", date, err)
	}
	defer rows.Close()
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

func (dao *BackupDaoSqlite) Delete(id string) error {
	_, err := dao.db.Exec("DELETE FROM backups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete backup => %s", err)
	}
	return nil
}
