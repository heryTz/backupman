package service

import (
	"fmt"
	"log"
	"os"

	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/model"
)

func HandleBackupStatus(app *core.App, id string) (model.Backup, error) {
	backup, err := app.Db.Backup.ReadBackupFullById(id)
	if err != nil {
		return model.Backup{}, fmt.Errorf("failed to read backup full by id (%s): %s", id, err)
	}
	if backup.Status == model.BACKUP_STATUS_FINISHED {
		return model.Backup{}, nil
	}

	countPending := 0
	countFailed := 0
	countFinished := 0

	for _, driveFile := range backup.DriveFiles {
		switch driveFile.Status {
		case model.DRIVE_FILE_STATUS_PENDING:
			countPending++
		case model.DRIVE_FILE_STATUS_FAILED:
			countFailed++
		case model.DRIVE_FILE_STATUS_FINISHED:
			countFinished++
		default:
			log.Fatalf("unknown drive file status (%s) for drive file (%s)", driveFile.Status, driveFile.Id)
		}
	}

	sampleBackup := model.Backup{
		Id:        backup.Id,
		Status:    backup.Status,
		Label:     backup.Label,
		DumpPath:  backup.DumpPath,
		CreatedAt: backup.CreatedAt,
	}

	if countPending > 0 {
		sampleBackup.Status = model.BACKUP_STATUS_PENDING
		_, err := app.Db.Backup.Update(backup.Id, sampleBackup)
		if err != nil {
			return model.Backup{}, fmt.Errorf("failed to update backup (%s) status to pending: %s", backup.Id, err)
		}
	} else if countFailed > 0 {
		sampleBackup.Status = model.BACKUP_STATUS_FAILED
		_, err := app.Db.Backup.Update(backup.Id, sampleBackup)
		if err != nil {
			return model.Backup{}, fmt.Errorf("failed to update backup (%s) status to failed: %s", backup.Id, err)
		}
	}

	sampleBackup.Status = model.BACKUP_STATUS_FINISHED
	_, err = app.Db.Backup.Update(backup.Id, sampleBackup)
	if err != nil {
		return model.Backup{}, fmt.Errorf("failed to update backup (%s) status to finished: %s", backup.Id, err)
	}

	return sampleBackup, nil
}

func RemoveBackupDump(app *core.App, backup model.Backup) error {
	if backup.DumpPath == "" {
		return nil
	}
	err := os.Remove(backup.DumpPath)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			log.Printf("dump file for backup %s does not exist => %s", backup.Id, err)
		default:
			return fmt.Errorf("cannot delete dump file %s => %s", backup.DumpPath, err)
		}
	}

	backup.DumpPath = ""
	_, err = app.Db.Backup.Update(backup.Id, backup)
	if err != nil {
		return err
	}

	return nil
}
