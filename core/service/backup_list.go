package service

import (
	"fmt"

	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/model"
)

type BackupListOutput struct {
	Results []model.BackupFull
}

func BackupList(app *core.App) (BackupListOutput, error) {
	data := BackupListOutput{
		Results: make([]model.BackupFull, 0),
	}
	results, err := app.Db.Backup.ReadAllBackupFull()
	if err != nil {
		return data, fmt.Errorf("error reading backups => %s", err)
	}
	data.Results = results
	return data, nil
}
