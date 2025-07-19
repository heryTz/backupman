//go:build test_integration

package tests_test

import (
	"testing"

	"github.com/herytz/backupman/core/drive"
	"github.com/stretchr/testify/assert"
)

var clientSecretFile = "../google-client-secret.json"
var tokenFile = "../google-token.json"

func TestGoogleDriveUploadFile(t *testing.T) {
	googleDrive := drive.NewGoogleDrive("Google Drive", "backupman", clientSecretFile, tokenFile)
	driveFile, err := googleDrive.Upload("./tmp/test.txt")
	assert.NoError(t, err)
	assert.NotEmpty(t, driveFile.Path)
}

func TestGoogleDriveDeleteFile(t *testing.T) {
	googleDrive := drive.NewGoogleDrive("Google Drive", "backupman", clientSecretFile, tokenFile)
	driveFile, err := googleDrive.Upload("./tmp/test.txt")
	assert.NoError(t, err)
	err = googleDrive.Delete(driveFile.Path)
	assert.NoError(t, err)
}
