package drive_test

import (
	"os"
	"testing"

	"github.com/herytz/backupman/core/drive"
)

func TestLocalDriveUploadFile(t *testing.T) {
	localDrive := drive.NewLocalDrive("local_drive", "./tmp/out")
	file, err := localDrive.Upload("./tmp/test.txt")
	if err != nil {
		t.Errorf("failed to upload file: %s", err)
	}
	_, err = os.Stat(file.Path)
	if err != nil {
		t.Errorf("file not upload correctly: %s", err)
	}
}
