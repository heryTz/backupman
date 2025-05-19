package model

import "time"

const (
	BACKUP_STATUS_PENDING  = "pending"
	BACKUP_STATUS_FINISHED = "finished"
	BACKUP_STATUS_FAILED   = "failed"
)

type Backup struct {
	Id        string
	Label     string
	Status    string
	DumpPath  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Backup) GetId() string {
	return b.Id
}

func (b *Backup) SetId(id string) {
	b.Id = id
}

type BackupFull struct {
	Id         string
	Status     string
	Label      string
	DumpPath   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DriveFiles []*DriveFile
}
