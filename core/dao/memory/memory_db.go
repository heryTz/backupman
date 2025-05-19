package memory

import "github.com/herytz/backupman/core/model"

type MemoryDb struct {
	Backup    *MemoryDbCrud[*model.Backup]
	DriveFile *MemoryDbCrud[*model.DriveFile]
}

func NewMemoryDb() *MemoryDb {
	return &MemoryDb{
		Backup:    NewMemoryDbCrud[*model.Backup]("backup"),
		DriveFile: NewMemoryDbCrud[*model.DriveFile]("drive_file"),
	}
}
