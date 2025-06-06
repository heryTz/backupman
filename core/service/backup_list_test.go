package service_test

import (
	"testing"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/service"
	"github.com/stretchr/testify/assert"
)

func TestBackupList(t *testing.T) {
	app := application.NewAppMock()
	backup1Ids, err := service.Backup(app)
	assert.NoError(t, err)
	backup2Ids, err := service.Backup(app)
	assert.NoError(t, err)
	list, err := service.BackupList(app)
	assert.NoError(t, err)
	var listIds []string
	for _, backup := range list.Results {
		listIds = append(listIds, backup.Id)
	}
	for _, backup1 := range backup1Ids {
		assert.Contains(t, listIds, backup1)
	}
	for _, backup2 := range backup2Ids {
		assert.Contains(t, listIds, backup2)
	}
}
