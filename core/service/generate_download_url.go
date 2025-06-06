package service

import (
	"fmt"
	"net/url"

	"github.com/herytz/backupman/core/application"
)

func GenerateDownloadUrl(app *application.App, backupId string) (string, error) {
	backup, err := app.Db.Backup.ReadFullById(backupId)
	if err != nil {
		return "", err
	}
	if backup == nil {
		return "", fmt.Errorf("backup with id %s not found", backupId)
	}

	// TODO: handle download preference from config

	downloadUrl := ""
	for _, driveFile := range backup.DriveFiles {
		if driveFile.Provider == "local" {
			downloadUrl, err = url.JoinPath(app.Http.AppUrl, fmt.Sprintf("api/backups/%s/download", driveFile.Id))
			if err != nil {
				return "", fmt.Errorf("failed to generate download url => %s", err)
			}
			break
		}
	}

	if downloadUrl == "" {
		return "", fmt.Errorf("no valid drive found for backup %s", backupId)
	}

	return downloadUrl, nil
}
