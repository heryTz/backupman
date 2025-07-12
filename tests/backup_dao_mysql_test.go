//go:build test_integration

package tests_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/herytz/backupman/core/dao/mysql"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/model"
	"github.com/stretchr/testify/assert"
)

var dbConn *sql.DB

func connectDb() {
	if dbConn == nil {
		conn, err := lib.NewConnection("localhost", 3307, "root", "root", "backupman", "false")
		if err != nil {
			log.Fatal(err)
		}
		dbConn = conn
	}
}

func TestReadFullById(t *testing.T) {
	connectDb()
	backupDao := mysql.NewBackupDaoMysql(dbConn)
	driveFileDao := mysql.NewDriveFileDaoMysql(dbConn)
	backupInput := model.Backup{
		Status:   model.BACKUP_STATUS_FINISHED,
		Label:    "backupLabel",
		DumpPath: "/tmp/backup.sql",
	}
	backup, err := backupDao.Create(backupInput)
	assert.NoError(t, err)
	driveFile1Input := model.DriveFile{
		BackupId: backup,
		Status:   model.DRIVE_FILE_STATUS_FINISHED,
		Path:     "/tmp/drive_file1",
		Label:    "driveLabel1",
		Provider: "local",
	}
	driveFile1, err := driveFileDao.Create(driveFile1Input)
	assert.NoError(t, err)
	driveFile2Input := model.DriveFile{
		BackupId: backup,
		Status:   model.DRIVE_FILE_STATUS_FINISHED,
		Path:     "/tmp/drive_file1",
		Label:    "driveLabel2",
		Provider: "local",
	}
	driveFile2, err := driveFileDao.Create(driveFile2Input)
	assert.NoError(t, err)

	backupFull, err := backupDao.ReadFullById(backup)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, backupFull)
	assert.Equal(t, backup, backupFull.Id)
	assert.Equal(t, backupInput.Status, backupFull.Status)
	assert.Equal(t, backupInput.Label, backupFull.Label)
	assert.Equal(t, backupInput.DumpPath, backupFull.DumpPath)

	for _, driveFile := range backupFull.DriveFiles {
		assert.Contains(t, []string{driveFile1, driveFile2}, driveFile.Id)
		assert.Contains(t, []string{driveFile1Input.Status, driveFile2Input.Status}, driveFile.Status)
		assert.Contains(t, []string{driveFile1Input.Path, driveFile2Input.Path}, driveFile.Path)
		assert.Contains(t, []string{driveFile1Input.Label, driveFile2Input.Label}, driveFile.Label)
		assert.Contains(t, []string{driveFile1Input.Provider, driveFile2Input.Provider}, driveFile.Provider)
	}
}
