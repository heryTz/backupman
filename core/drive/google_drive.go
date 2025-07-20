package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gdrive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDrive struct {
	Label            string
	Folder           string
	ClientSecretFile string
	TokenFile        string
}

func NewGoogleDrive(label, folder, clientSecretFile, tokenFile string) *GoogleDrive {
	drive := GoogleDrive{
		Label:            label,
		Folder:           folder,
		ClientSecretFile: clientSecretFile,
		TokenFile:        tokenFile,
	}
	return &drive
}

func (d *GoogleDrive) getDriveService() (*gdrive.Service, error) {
	ctx := context.Background()
	b, err := os.ReadFile(d.ClientSecretFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gdrive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	f, err := os.Open(d.TokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open token file: %v", err)
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, fmt.Errorf("unable to decode token file: %v", err)
	}

	client := config.Client(ctx, tok)
	return gdrive.NewService(ctx, option.WithHTTPClient(client))
}

func (d *GoogleDrive) findOrCreateFolder(srv *gdrive.Service) (*gdrive.File, error) {
	query := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and trashed=false", d.Folder)
	files, err := srv.Files.List().
		Q(query).
		Spaces("drive").
		Fields("files(id, name)").
		Do()
	if err != nil {
		return nil, fmt.Errorf("error retrieving folder => %v", err)
	}

	if len(files.Files) > 0 {
		return files.Files[0], nil
	}

	folderMetadata := &gdrive.File{
		Name:     d.Folder,
		MimeType: "application/vnd.google-apps.folder",
	}
	folder, err := srv.Files.Create(folderMetadata).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create folder => %v", err)
	}

	return folder, nil
}

func (d *GoogleDrive) Upload(srcPath string) (DriveFile, error) {
	driveFile := DriveFile{}
	srv, err := d.getDriveService()
	if err != nil {
		return driveFile, fmt.Errorf("[Google Drive] Unable to retrieve Drive client => %s", err)
	}

	file, err := os.Open(srcPath)
	if err != nil {
		return driveFile, fmt.Errorf("[Google Drive] Unable to open file %s => %s", srcPath, err)
	}
	defer file.Close()

	folder, err := d.findOrCreateFolder(srv)
	if err != nil {
		return driveFile, err
	}

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
	srv, err := d.getDriveService()
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

func (d *GoogleDrive) Health() error {
	folder := "./tmp"
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create temporary directory for healh test => %s", err)
	}
	healthTest := path.Join(folder, "health_test.txt")
	os.Remove(healthTest)
	err = os.WriteFile(healthTest, []byte("health test"), 0755)
	if err != nil {
		return fmt.Errorf("Failed to create health test file => %s", err)
	}

	file, err := d.Upload(healthTest)
	if err != nil {
		return fmt.Errorf("Failed to upload health test file to Google Drive => %s", err)
	}

	err = d.Delete(file.Path)
	if err != nil {
		return fmt.Errorf("Failed to delete health test file from Google Drive => %s", err)
	}

	return nil
}

func (d *GoogleDrive) GetLabel() string {
	return d.Label
}

func (d *GoogleDrive) GetProvider() string {
	return "google_drive"
}
