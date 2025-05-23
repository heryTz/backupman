package drive

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LocalDrive struct {
	Label  string
	Folder string
}

func NewLocalDrive(label, folder string) *LocalDrive {
	drive := LocalDrive{
		Label:  label,
		Folder: folder,
	}
	drive.setup()
	return &drive
}

func (d *LocalDrive) Upload(srcPath string) (DriveFile, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Printf("failed to open file => %s", err)
		return DriveFile{}, err
	}
	defer srcFile.Close()

	dstFilename := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), filepath.Ext(srcPath))
	dstPath := filepath.Join(d.Folder, dstFilename)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Printf("failed to create file => %s", err)
		return DriveFile{}, err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	if err != nil {
		log.Printf("failed to copy file => %s", err)
		return DriveFile{}, err
	}

	return DriveFile{
		Path: dstFile.Name(),
	}, nil
}

func (d *LocalDrive) Delete(srcPath string) error {
	err := os.Remove(srcPath)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			log.Printf("file does not exist => %s", err)
		default:
			log.Printf("failed to delete file: %s", err)
			return err
		}
	}
	return nil
}

func (d *LocalDrive) GetLabel() string {
	return d.Label
}

func (d *LocalDrive) setup() {
	err := os.MkdirAll(d.Folder, 0755)
	if err != nil {
		log.Fatalf("failed to setup local drive. %s", err)
	}
}
func (d *LocalDrive) GetProvider() string {
	return "local"
}
