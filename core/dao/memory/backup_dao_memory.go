package memory

import (
	"fmt"

	"github.com/herytz/backupman/core/model"
)

type BackupDaoMemory struct {
	db *MemoryDb
}

func NewBackupDaoMemory(db *MemoryDb) *BackupDaoMemory {
	return &BackupDaoMemory{
		db: db,
	}
}

func (dao *BackupDaoMemory) ReadOrError(id string) (*model.Backup, error) {
	backup := dao.db.Backup.ReadById(id)
	if backup == nil {
		return nil, fmt.Errorf("no backup found with id %s", id)
	}
	return backup, nil
}

func (dao *BackupDaoMemory) Create(data model.Backup) (string, error) {
	result, err := dao.db.Backup.Create(&data)
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (dao *BackupDaoMemory) Update(id string, data model.Backup) (string, error) {
	_, err := dao.db.Backup.Update(id, &data)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (dao *BackupDaoMemory) ReadBackupFullById(id string) (*model.BackupFull, error) {
	backup := dao.db.Backup.ReadById(id)
	if backup == nil {
		return nil, nil
	}
	driveFiles := dao.db.DriveFile.ReadAll()
	var backupDriveFiles []*model.DriveFile
	for _, driveFile := range driveFiles {
		if driveFile.BackupId == id {
			backupDriveFiles = append(backupDriveFiles, driveFile)
		}
	}
	backupFull := &model.BackupFull{
		Id:         backup.Id,
		Status:     backup.Status,
		Label:      backup.Label,
		DumpPath:   backup.DumpPath,
		CreatedAt:  backup.CreatedAt,
		DriveFiles: backupDriveFiles,
	}
	return backupFull, nil
}

func (dao *BackupDaoMemory) ReadAllBackupFull() ([]model.BackupFull, error) {
	backupFullList := make([]model.BackupFull, 0)
	backupList := dao.db.Backup.ReadAll()
	driveFiles := dao.db.DriveFile.ReadAll()
	for _, backup := range backupList {
		var backupDriveFiles []*model.DriveFile
		for _, driveFile := range driveFiles {
			if driveFile.BackupId == backup.Id {
				backupDriveFiles = append(backupDriveFiles, driveFile)
			}
		}
		backupFull := model.BackupFull{
			Id:         backup.Id,
			Status:     backup.Status,
			Label:      backup.Label,
			DumpPath:   backup.DumpPath,
			CreatedAt:  backup.CreatedAt,
			DriveFiles: backupDriveFiles,
		}
		backupFullList = append(backupFullList, backupFull)
	}
	return backupFullList, nil
}
