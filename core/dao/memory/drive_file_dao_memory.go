package memory

import (
	"fmt"

	"github.com/herytz/backupman/core/model"
)

type DriveFileDaoMemory struct {
	db *MemoryDb
}

func NewDriveFileDaoMemory(db *MemoryDb) *DriveFileDaoMemory {
	return &DriveFileDaoMemory{
		db: db,
	}
}

func (dao *DriveFileDaoMemory) Create(data model.DriveFile) (string, error) {
	result, err := dao.db.DriveFile.Create(&data)
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (dao *DriveFileDaoMemory) Update(id string, data model.DriveFile) (string, error) {
	result, err := dao.db.DriveFile.Update(id, &data)
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (dao *DriveFileDaoMemory) ReadOrError(id string) (*model.DriveFile, error) {
	driveFile := dao.db.DriveFile.ReadById(id)
	if driveFile == nil {
		return nil, fmt.Errorf("no drive file found with id %s", id)
	}
	return driveFile, nil
}

func (dao *DriveFileDaoMemory) ReadById(id string) *model.DriveFile {
	return dao.db.DriveFile.ReadById(id)
}

func (dao *DriveFileDaoMemory) ReadAll() []*model.DriveFile {
	return dao.db.DriveFile.ReadAll()
}
