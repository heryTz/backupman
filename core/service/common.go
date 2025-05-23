package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/drive"
	"github.com/herytz/backupman/core/model"
)

func HandleBackupStatus(app *core.App, id string) (model.Backup, error) {
	backup, err := app.Db.Backup.ReadFullById(id)
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

func RemoveOldBackup(app *core.App) error {
	maxAge := time.Now().AddDate(0, 0, -app.Retention.Days)
	backups, err := app.Db.Backup.ReadOlderThan(maxAge)
	if err != nil {
		return fmt.Errorf("failed to read backups older than %s => %s", maxAge, err)
	}

	for _, backup := range backups {
		for _, driveFile := range backup.DriveFiles {
			var drive drive.Drive
			for _, d := range app.Drives {
				if d.GetProvider() == driveFile.Provider {
					drive = d
					break
				}
			}
			if drive == nil {
				return fmt.Errorf("drive not found for provider %s", driveFile.Provider)
			}
			err := drive.Delete(driveFile.Path)
			if err != nil {
				return fmt.Errorf("failed to delete drive file (%s) => %s", driveFile.Path, err)
			}
			err = app.Db.DriveFile.Delete(driveFile.Id)
			if err != nil {
				return fmt.Errorf("failed to delete drive file (%s) => %s", driveFile.Id, err)
			}
		}

		err := app.Db.Backup.Delete(backup.Id)
		if err != nil {
			return fmt.Errorf("failed to delete backup (%s) => %s", backup.Id, err)
		}
	}

	return nil
}
