package dao

import (
	"time"

	"github.com/herytz/backupman/core/model"
)

type BackupDao interface {
	Create(data model.Backup) (string, error)
	Update(id string, data model.Backup) (string, error)
	ReadFullById(Id string) (*model.BackupFull, error)
	ReadAllFull() ([]model.BackupFull, error)
	ReadOrError(id string) (*model.Backup, error)
	ReadOlderThan(date time.Time) ([]model.BackupFull, error)
	Delete(id string) error
}

type DriveFileDao interface {
	Create(data model.DriveFile) (string, error)
	Update(Id string, data model.DriveFile) (string, error)
	ReadOrError(id string) (*model.DriveFile, error)
	Delete(id string) error
}

type Dao struct {
	Backup    BackupDao
	DriveFile DriveFileDao
}
