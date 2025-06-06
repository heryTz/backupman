//go:build test_integration

package drive_test

import (
	"testing"

	"github.com/herytz/backupman/core/drive"
	"github.com/stretchr/testify/assert"
)

func TestGoogleDriveUploadFile(t *testing.T) {
	serviceAccount := "./service-account.json"
	googleDrive := drive.NewGoogleDrive("Google Drive", "demo", serviceAccount)
	driveFile, err := googleDrive.Upload("./tmp/test.txt")
	assert.NoError(t, err)
	assert.NotEmpty(t, driveFile.Path)
}

func TestGoogleDriveDeleteFile(t *testing.T) {
	serviceAccount := "./service-account.json"
	googleDrive := drive.NewGoogleDrive("Google Drive", "demo", serviceAccount)
	driveFile, err := googleDrive.Upload("./tmp/test.txt")
	assert.NoError(t, err)
	err = googleDrive.Delete(driveFile.Path)
	assert.NoError(t, err)
}
