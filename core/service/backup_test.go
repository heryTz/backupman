package service_test

import (
	"testing"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/model"
	"github.com/herytz/backupman/core/service"
	"github.com/stretchr/testify/assert"
)

func TestBackupReport(t *testing.T) {
	app := application.NewAppMock()
	backupIds, err := service.Backup(app)
	assert.NoError(t, err)
	for _, backupId := range backupIds {
		assert.NotEqual(t, "", backupId)
		backupFull, err := app.Db.Backup.ReadFullById(backupId)
		assert.NoError(t, err)
		assert.NotEqual(t, nil, backupFull)
		assert.Equal(t, model.BACKUP_STATUS_FINISHED, backupFull.Status)
		assert.NotEqual(t, 0, len(backupFull.DriveFiles))
		driveFileFinished := 0
		for _, driveFile := range backupFull.DriveFiles {
			if driveFile.Status == model.DRIVE_FILE_STATUS_FINISHED {
				driveFileFinished++
			}
		}
		assert.Equal(t, len(backupFull.DriveFiles), driveFileFinished)
	}
}
