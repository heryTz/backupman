//go:build test_integration

package tests_test

import (
	"testing"

	"github.com/herytz/backupman/core/drive"
	"github.com/stretchr/testify/assert"
)

func TestLocalDriveUploadFile(t *testing.T) {
	localDrive := drive.NewLocalDrive("local_drive", "./tmp/out")
	file, err := localDrive.Upload("./tmp/test.txt")
	assert.NoError(t, err)
	assert.FileExists(t, file.Path)
}

func TestLocalDriveDeleteFile(t *testing.T) {
	localDrive := drive.NewLocalDrive("local_drive", "./tmp/out")
	file, err := localDrive.Upload("./tmp/test.txt")
	assert.NoError(t, err)
	err = localDrive.Delete(file.Path)
	assert.NoError(t, err)
	assert.NoFileExists(t, file.Path)
}
