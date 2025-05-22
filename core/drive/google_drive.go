package drive

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	gdrive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDrive struct {
	Label          string
	Folder         string
	serviceAccount string
}

func NewGoogleDrive(label, folder, serviceAccount string) *GoogleDrive {
	drive := GoogleDrive{
		Label:          label,
		Folder:         folder,
		serviceAccount: serviceAccount,
	}
	return &drive
}

func (d *GoogleDrive) Upload(srcPath string) (DriveFile, error) {
	driveFile := DriveFile{}
	ctx := context.Background()

	srv, err := gdrive.NewService(ctx, option.WithCredentialsFile(d.serviceAccount))
	if err != nil {
		return driveFile, fmt.Errorf("[Google Drive] Unable to retrieve Drive client => %s", err)
	}

	file, err := os.Open(srcPath)
	if err != nil {
		return driveFile, fmt.Errorf("[Google Drive] Unable to open file %s => %s", srcPath, err)
	}
	defer file.Close()

	query := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s'", d.Folder)
	files, err := srv.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des dossiers: %v", err)
	}

	if len(files.Files) == 0 {
		return driveFile, fmt.Errorf("[Google Drive] Folder %s not found", d.Folder)
	}

	folder := files.Files[0]
	filename := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), filepath.Ext(srcPath))

	fileMetadata := &gdrive.File{
		Name:    filename,
		Parents: []string{folder.Id},
	}

	uploadedFile, err := srv.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		return driveFile, fmt.Errorf("[Google Drive] Unable to upload file %s => %s", srcPath, err)
	}

	driveFile.Path = uploadedFile.Name

	return driveFile, nil
}

func (d *GoogleDrive) Delete(srcPath string) error {
	ctx := context.Background()

	srv, err := gdrive.NewService(ctx, option.WithCredentialsFile(d.serviceAccount))
	if err != nil {
		return fmt.Errorf("[Google Drive] Unable to retrieve Drive client => %s", err)
	}

	query := fmt.Sprintf("name='%s'", filepath.Base(srcPath))
	files, err := srv.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des fichiers => %v", err)
	}

	if len(files.Files) == 0 {
		log.Printf("[Google Drive] File %s not found", srcPath)
		return nil
	}

	for _, file := range files.Files {
		err = srv.Files.Delete(file.Id).Do()
		if err != nil {
			return fmt.Errorf("[Google Drive] Unable to delete file %s => %s", srcPath, err)
		}
	}

	return nil
}

func (d *GoogleDrive) GetLabel() string {
	return d.Label
}

func (d *GoogleDrive) GetProvider() string {
	return "google_drive"
}
