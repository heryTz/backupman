//go:build test_integration

package tests_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/dao/sqlite"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/model"
	"github.com/herytz/backupman/migration"
	"github.com/stretchr/testify/assert"
)

var sqliteDbConn *sql.DB

func connectSqliteDb() {
	dbPath := "/tmp/backupman_test.db"

	os.Remove(dbPath)

	conn, err := lib.NewSqliteConnection(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	sqliteDbConn = conn

	err = migration.Run(application.SqliteDbConfig{
		DbPath: dbPath,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestSqliteReadFullById(t *testing.T) {
	connectSqliteDb()
	defer sqliteDbConn.Close()

	backupDao := sqlite.NewBackupDaoSqlite(sqliteDbConn)
	driveFileDao := sqlite.NewDriveFileDaoSqlite(sqliteDbConn)

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
		Path:     "/tmp/drive_file2",
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

func TestSqliteReadAllFull(t *testing.T) {
	connectSqliteDb()
	defer sqliteDbConn.Close()

	backupDao := sqlite.NewBackupDaoSqlite(sqliteDbConn)
	driveFileDao := sqlite.NewDriveFileDaoSqlite(sqliteDbConn)

	backupInput := model.Backup{
		Status:   model.BACKUP_STATUS_FINISHED,
		Label:    "backupLabel",
		DumpPath: "/tmp/backup.sql",
	}
	backup, err := backupDao.Create(backupInput)
	assert.NoError(t, err)

	driveFileInput := model.DriveFile{
		BackupId: backup,
		Status:   model.DRIVE_FILE_STATUS_FINISHED,
		Path:     "/tmp/drive_file",
		Label:    "driveLabel",
		Provider: "local",
	}
	_, err = driveFileDao.Create(driveFileInput)
	assert.NoError(t, err)

	backups, err := backupDao.ReadAllFull()
	assert.NoError(t, err)
	assert.Greater(t, len(backups), 0)
}

func TestSqliteCreateAndUpdate(t *testing.T) {
	connectSqliteDb()
	defer sqliteDbConn.Close()

	backupDao := sqlite.NewBackupDaoSqlite(sqliteDbConn)

	backupInput := model.Backup{
		Status:   model.BACKUP_STATUS_PENDING,
		Label:    "backupLabel",
		DumpPath: "/tmp/backup.sql",
	}
	backupId, err := backupDao.Create(backupInput)
	assert.NoError(t, err)
	assert.NotEmpty(t, backupId)

	// Update backup
	updateInput := model.Backup{
		Status:   model.BACKUP_STATUS_FINISHED,
		Label:    "updatedLabel",
		DumpPath: "/tmp/backup_updated.sql",
	}
	_, err = backupDao.Update(backupId, updateInput)
	assert.NoError(t, err)

	// Read and verify
	backup, err := backupDao.ReadOrError(backupId)
	assert.NoError(t, err)
	assert.Equal(t, updateInput.Status, backup.Status)
	assert.Equal(t, updateInput.Label, backup.Label)
	assert.Equal(t, updateInput.DumpPath, backup.DumpPath)
}

func TestSqliteDelete(t *testing.T) {
	connectSqliteDb()
	defer sqliteDbConn.Close()

	backupDao := sqlite.NewBackupDaoSqlite(sqliteDbConn)

	backupInput := model.Backup{
		Status:   model.BACKUP_STATUS_FINISHED,
		Label:    "backupLabel",
		DumpPath: "/tmp/backup.sql",
	}
	backupId, err := backupDao.Create(backupInput)
	assert.NoError(t, err)

	// Delete backup
	err = backupDao.Delete(backupId)
	assert.NoError(t, err)

	// Verify deletion
	backup, err := backupDao.ReadOrError(backupId)
	assert.Error(t, err)
	assert.Nil(t, backup)
}
