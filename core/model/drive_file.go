package model

import "time"

const (
	DRIVE_FILE_STATUS_PENDING  = "pending"
	DRIVE_FILE_STATUS_FINISHED = "finished"
	DRIVE_FILE_STATUS_FAILED   = "failed"
)

type DriveFile struct {
	Id        string
	BackupId  string
	Provider  string
	Label     string
	Path      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *DriveFile) GetId() string {
	return b.Id
}

func (b *DriveFile) SetId(id string) {
	b.Id = id
}
