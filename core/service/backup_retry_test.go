package service_test

import (
	"os"
	"path"
	"testing"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/model"
	"github.com/herytz/backupman/core/service"
	"github.com/stretchr/testify/assert"
)

func TestBackupRetry(t *testing.T) {
	app := application.NewAppMock()
	backupIds, err := service.Backup(app)
	assert.NoError(t, err)
	assert.NotEmpty(t, backupIds)
	backupId := backupIds[0]

	backup, err := app.Db.Backup.ReadOrError(backupId)
	assert.NoError(t, err)
	backup.Status = model.BACKUP_STATUS_FAILED
	dummyFolder := "./tmp"
	err = os.MkdirAll(dummyFolder, 0755)
	assert.NoError(t, err, "Failed to create temporary directory for backup")
	dummyBackup := path.Join(dummyFolder, "dummy_backup.sql")
	_, err = os.Create(dummyBackup)
	assert.NoError(t, err, "Failed to create dummy backup file")
	backup.DumpPath = dummyBackup
	_, err = app.Db.Backup.Update(backupId, *backup)
	assert.NoError(t, err)

	backupFull, err := app.Db.Backup.ReadFullById(backupId)
	assert.NoError(t, err)
	for _, driveFile := range backupFull.DriveFiles {
		driveFile.Status = model.DRIVE_FILE_STATUS_FAILED
		_, err = app.Db.DriveFile.Update(driveFile.Id, *driveFile)
		assert.NoError(t, err)
	}

	err = service.BackupRetry(app, backupId)
	assert.NoError(t, err)
	newBackup, err := app.Db.Backup.ReadOrError(backupId)
	assert.NoError(t, err)
	assert.Empty(t, newBackup.DumpPath)
	assert.Equal(t, model.BACKUP_STATUS_FINISHED, newBackup.Status)
}
