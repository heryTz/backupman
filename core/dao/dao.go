package dao

import "github.com/herytz/backupman/core/model"

type BackupDao interface {
	Create(data model.Backup) (string, error)
	Update(id string, data model.Backup) (string, error)
	ReadBackupFullById(Id string) (*model.BackupFull, error)
	ReadAllBackupFull() ([]model.BackupFull, error)
	ReadOrError(id string) (*model.Backup, error)
}

type DriveFileDao interface {
	Create(data model.DriveFile) (string, error)
	Update(Id string, data model.DriveFile) (string, error)
	ReadOrError(id string) (*model.DriveFile, error)
}

type Dao struct {
	Backup    BackupDao
	DriveFile DriveFileDao
}
