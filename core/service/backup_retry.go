package service

import (
	"fmt"
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/drive"
	"github.com/herytz/backupman/core/model"
)

func BackupRetry(app *application.App, backupId string) error {
	backup, err := app.Db.Backup.ReadFullById(backupId)
	if err != nil {
		return fmt.Errorf("failed to read backup => %s", err)
	}
	if backup.DumpPath == "" {
		return fmt.Errorf("backup dump path is empty")
	}
	if backup.Status != model.BACKUP_STATUS_FAILED {
		return fmt.Errorf("backup status is not failed")
	}

	for _, driveFile := range backup.DriveFiles {
		if driveFile.Status != model.DRIVE_FILE_STATUS_FAILED {
			continue
		}

		uploadResult, err := upload(app, backup.DumpPath, driveFile)
		if err != nil {
			log.Printf("failed to upload dump (%s) for database (%s) to drive (%s) => %s", backup.DumpPath, backup.Label, driveFile.Provider, err)

			driveFile.Status = model.DRIVE_FILE_STATUS_FAILED
			_, err := app.Db.DriveFile.Update(driveFile.Id, *driveFile)
			if err != nil {
				log.Printf("failed to update drive file (%s) status to failed => %s", driveFile.Id, err)
			}
			continue
		}

		driveFile.Status = model.DRIVE_FILE_STATUS_FINISHED
		driveFile.Path = uploadResult.Path
		_, err = app.Db.DriveFile.Update(driveFile.Id, *driveFile)
		if err != nil {
			log.Printf("failed to update drive file (%s) status to finished => %s", driveFile.Id, err)
			continue
		}
	}

	err = AfterBackup(app, backupId)
	if err != nil {
		return fmt.Errorf("failed to execute after backup actions => %s", err)
	}

	return nil
}

func upload(app *application.App, dumpPath string, driveFile *model.DriveFile) (drive.DriveFile, error) {
	var uploadResult drive.DriveFile

	driveFile.Status = model.DRIVE_FILE_STATUS_PENDING
	_, err := app.Db.DriveFile.Update(driveFile.Id, *driveFile)
	if err != nil {
		return uploadResult, fmt.Errorf("failed to update drive file (%s) status to pending => %s", driveFile.Id, err)
	}

	drive, err := GetDrive(app, driveFile.Provider)
	if err != nil {
		return uploadResult, fmt.Errorf("failed to get drive (%s) => %s", driveFile.Provider, err)
	}

	uploadResult, err = drive.Upload(dumpPath)
	if err != nil {
		return uploadResult, fmt.Errorf("failed to upload dump (%s) => %s", dumpPath, err)
	}

	return uploadResult, nil
}
