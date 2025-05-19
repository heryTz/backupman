package service

import (
	"fmt"
	"os"

	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/model"
)

type DownloadOutput struct {
	Filename string
	Byte     []byte
	MimeType string
}

func Download(app *core.App, driveFileId string) (DownloadOutput, error) {
	var output DownloadOutput
	driveFile, err := app.Db.DriveFile.ReadOrError(driveFileId)
	if err != nil {
		return output, err
	}

	if driveFile.Status != model.DRIVE_FILE_STATUS_FINISHED {
		return output, fmt.Errorf("drive file %s is not finished", driveFileId)
	}

	if driveFile.Provider != "local" {
		return output, fmt.Errorf("unexpected drive provider (%s). This route only supports local drive", driveFile.Provider)
	}

	data, err := os.ReadFile(driveFile.Path)
	if err != nil {
		return output, fmt.Errorf("failed to read file %s => %s", driveFile.Path, err)
	}

	output.Byte = data
	output.Filename = fmt.Sprintf("%s-%s.sql", driveFile.Label, driveFile.CreatedAt.Format("2006-01-02_15-04-05"))
	output.MimeType = "application/octet-stream"

	return output, nil
}
