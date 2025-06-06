package service

import (
	"fmt"
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/model"
)

func Backup(app *application.App) ([]string, error) {
	backupIds := make([]string, 0)

	for _, dumper := range app.Dumpers {
		backupId, err := app.Db.Backup.Create(model.Backup{
			Label:  dumper.GetLabel(),
			Status: model.BACKUP_STATUS_PENDING,
		})
		if err != nil {
			return backupIds, fmt.Errorf("failed to create backup => %s", err)
		}
		backupIds = append(backupIds, backupId)

		backup, err := app.Db.Backup.ReadOrError(backupId)
		if err != nil {
			return backupIds, fmt.Errorf("failed to read backup => %s", err)
		}

		dump, err := dumper.Dump()
		if err != nil {
			log.Printf("failed to dump database (%s) => %s", dumper.GetLabel(), err)
			backup.Status = model.BACKUP_STATUS_FAILED
			_, err := app.Db.Backup.Update(backup.Id, *backup)
			if err != nil {
				log.Printf("failed to update backup (%s) status to failed => %s", backup.Id, err)
			}

			continue
		}

		backup.DumpPath = dump
		_, err = app.Db.Backup.Update(backup.Id, *backup)
		if err != nil {
			log.Printf("failed to update backup (%s) with dump path (%s): %s", backup.Id, dump, err)
			continue
		}

		for _, drive := range app.Drives {
			driveFileId, err := app.Db.DriveFile.Create(model.DriveFile{
				BackupId: backup.Id,
				Status:   model.DRIVE_FILE_STATUS_PENDING,
				Label:    drive.GetLabel(),
				Provider: drive.GetProvider(),
			})
			if err != nil {
				log.Printf("failed to create drive (%s) for database (%s) => %s", drive.GetLabel(), dumper.GetLabel(), err)
				continue
			}

			driveFile, err := app.Db.DriveFile.ReadOrError(driveFileId)
			if err != nil {
				log.Printf("failed to read drive file (%s) => %s", driveFileId, err)
				continue
			}
			file, err := drive.Upload(dump)
			if err != nil {
				log.Printf("failed to upload dump (%s) for database (%s) to drive (%s) => %s", dump, dumper.GetLabel(), drive.GetLabel(), err)
				driveFile.Status = model.DRIVE_FILE_STATUS_FAILED
				_, err := app.Db.DriveFile.Update(driveFile.Id, *driveFile)
				if err != nil {
					log.Printf("failed to update drive file (%s) status to failed => %s", driveFile.Id, err)
				}
				continue
			}

			driveFile.Status = model.DRIVE_FILE_STATUS_FINISHED
			driveFile.Path = file.Path
			_, err = app.Db.DriveFile.Update(driveFile.Id, *driveFile)
			if err != nil {
				log.Printf("failed to update drive file (%s) status to finished => %s", driveFile.Id, err)
			}
		}

		err = AfterBackup(app, backupId)
		if err != nil {
			log.Printf("failed to execute after backup tasks for backup (%s) => %s", backupId, err)
		}
	}

	if app.Retention.Enabled {
		if app.Mode == application.APP_MODE_CLI {
			err := RemoveOldBackup(app)
			if err != nil {
				log.Println(err)
			}
		} else {
			go func(app *application.App) {
				err := RemoveOldBackup(app)
				if err != nil {
					log.Println(err)
				}
			}(app)
		}
	}

	return backupIds, nil
}
