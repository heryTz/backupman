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

func (d *GoogleDrive) getDriveService() (*gdrive.Service, error) {
	ctx := context.Background()
	switch d.AuthType {
	case "serviceAccount":
		return gdrive.NewService(ctx, option.WithCredentialsFile(d.ServiceAccountFile))
	case "oauth2":
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
	default:
		return nil, fmt.Errorf("unknown auth type: %s", d.AuthType)
	}
}

type GoogleDrive struct {
	Label              string
	Folder             string
	AuthType           string
	ServiceAccountFile string
	ClientSecretFile   string
	TokenFile          string
}

func NewGoogleDriveWithServiceAccount(label, folder, serviceAccount string) *GoogleDrive {
	drive := GoogleDrive{
		Label:              label,
		Folder:             folder,
		AuthType:           "serviceAccount",
		ServiceAccountFile: serviceAccount,
	}
	return &drive
}

func NewGoogleDriveWithOauth2(label, folder, clientSecretFile, tokenFile string) *GoogleDrive {
	drive := GoogleDrive{
		Label:            label,
		Folder:           folder,
		AuthType:         "oauth2",
		ClientSecretFile: clientSecretFile,
		TokenFile:        tokenFile,
	}
	return &drive
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
