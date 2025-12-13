package drive

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Drive struct {
	Label                string
	Bucket               string
	Region               string
	AccessKey            string
	SecretKey            string
	Endpoint             string
	Prefix               string
	ForcePathStyle       bool
	EnableIntegrityCheck bool
}

func NewS3Drive(label, bucket, region, accessKey, secretKey, endpoint, prefix string, forcePathStyle bool) *S3Drive {
	drive := S3Drive{
		Label:                label,
		Bucket:               bucket,
		Region:               region,
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		Endpoint:             endpoint,
		Prefix:               prefix,
		ForcePathStyle:       forcePathStyle,
		EnableIntegrityCheck: true,
	}
	return &drive
}

func (d *S3Drive) getS3Client() (*s3.Client, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(d.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	// Override credentials if provided
	if d.AccessKey != "" && d.SecretKey != "" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider(d.AccessKey, d.SecretKey, "")
	}

	// Create S3 client with custom endpoint if provided (for S3-compatible services)
	clientOptions := []func(*s3.Options){}
	if d.Endpoint != "" {
		clientOptions = append(clientOptions, func(o *s3.Options) {
			o.BaseEndpoint = &d.Endpoint
			o.UsePathStyle = d.ForcePathStyle
		})
	}

	return s3.NewFromConfig(cfg, clientOptions...), nil
}

func (d *S3Drive) Upload(srcPath string) (DriveFile, error) {
	driveFile := DriveFile{}

	client, err := d.getS3Client()
	if err != nil {
		return driveFile, fmt.Errorf("[S3 Drive] Unable to create S3 client => %s", err)
	}

	file, err := os.Open(srcPath)
	if err != nil {
		return driveFile, fmt.Errorf("[S3 Drive] Unable to open file %s => %s", srcPath, err)
	}
	defer file.Close()

	// Calculate file checksums for integrity verification
	var localMD5, localSHA256 string
	if d.EnableIntegrityCheck {
		localMD5, localSHA256, err = d.calculateChecksums(file)
		if err != nil {
			return driveFile, fmt.Errorf("[S3 Drive] Failed to calculate checksums => %s", err)
		}

		// Reset file pointer for upload
		_, err = file.Seek(0, 0)
		if err != nil {
			return driveFile, fmt.Errorf("[S3 Drive] Failed to reset file pointer => %s", err)
		}
	}

	filename := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), filepath.Ext(srcPath))

	// Build S3 key with prefix
	key := filename
	if d.Prefix != "" {
		key = fmt.Sprintf("%s/%s", strings.TrimPrefix(d.Prefix, "/"), filename)
	}

	// Upload with metadata for integrity verification
	putObjectInput := &s3.PutObjectInput{
		Bucket: &d.Bucket,
		Key:    &key,
		Body:   file,
	}

	if d.EnableIntegrityCheck {
		putObjectInput.Metadata = map[string]string{
			"local-md5":    localMD5,
			"local-sha256": localSHA256,
		}
	}

	result, err := client.PutObject(context.Background(), putObjectInput)
	if err != nil {
		return driveFile, fmt.Errorf("[S3 Drive] Unable to upload file %s => %s", srcPath, err)
	}

	// Verify upload integrity if enabled
	if d.EnableIntegrityCheck && result.ETag != nil {
		err = d.verifyUploadIntegrity(client, key, *result.ETag, localMD5, localSHA256)
		if err != nil {
			// Attempt to clean up failed upload
			_ = d.Delete(key)
			return driveFile, fmt.Errorf("[S3 Drive] Upload integrity verification failed => %s", err)
		}
	}

	driveFile.Path = key
	if result.ETag != nil {
		driveFile.Checksum = *result.ETag
	}

	return driveFile, nil
}

func (d *S3Drive) Delete(srcPath string) error {
	client, err := d.getS3Client()
	if err != nil {
		return fmt.Errorf("[S3 Drive] Unable to create S3 client => %s", err)
	}

	// Use the full path as key since Upload returns the full S3 key
	key := srcPath

	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &d.Bucket,
		Key:    &key,
	})

	if err != nil {
		return fmt.Errorf("[S3 Drive] Unable to delete file %s => %s", srcPath, err)
	}

	return nil
}

func (d *S3Drive) Health() error {
	folder := "./tmp"
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create temporary directory for health test => %s", err)
	}

	healthTest := filepath.Join(folder, "health_test.txt")
	os.Remove(healthTest)
	err = os.WriteFile(healthTest, []byte("health test"), 0755)
	if err != nil {
		return fmt.Errorf("Failed to create health test file => %s", err)
	}

	file, err := d.Upload(healthTest)
	if err != nil {
		return fmt.Errorf("Failed to upload health test file to S3 => %s", err)
	}

	err = d.Delete(file.Path)
	if err != nil {
		return fmt.Errorf("Failed to delete health test file from S3 => %s", err)
	}

	// Clean up local file
	os.Remove(healthTest)

	return nil
}

func (d *S3Drive) GetLabel() string {
	return d.Label
}

func (d *S3Drive) GetProvider() string {
	return "s3"
}

// calculateChecksums computes MD5 and SHA256 hashes for the file
func (d *S3Drive) calculateChecksums(file io.ReadSeeker) (string, string, error) {
	md5Hash := md5.New()
	sha256Hash := sha256.New()

	// Create a multi-writer to calculate both hashes simultaneously
	multiWriter := io.MultiWriter(md5Hash, sha256Hash)

	// Copy file data through the multi-writer
	_, err := io.Copy(multiWriter, file)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file for checksum calculation: %w", err)
	}

	md5Sum := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Sum := hex.EncodeToString(sha256Hash.Sum(nil))

	return md5Sum, sha256Sum, nil
}

// verifyUploadIntegrity verifies the uploaded file integrity using S3 ETag
func (d *S3Drive) verifyUploadIntegrity(client *s3.Client, key, etag, localMD5, localSHA256 string) error {
	// Get object metadata to verify ETag
	headResult, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &d.Bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to get object metadata: %w", err)
	}

	// Verify ETag matches
	if headResult.ETag == nil {
		return fmt.Errorf("no ETag returned from S3")
	}

	// Clean ETag (remove quotes if present)
	s3ETag := strings.Trim(*headResult.ETag, "\"")
	localMD5Clean := strings.Trim(localMD5, "\"")

	// For single-part uploads, ETag should be MD5 hash
	// For multi-part uploads, ETag format is different
	if !strings.Contains(s3ETag, "-") {
		// Single-part upload - compare MD5
		if s3ETag != localMD5Clean {
			return fmt.Errorf("ETag mismatch: S3 ETag=%s, local MD5=%s", s3ETag, localMD5Clean)
		}
	}

	// Verify metadata was stored correctly
	if headResult.Metadata == nil {
		return fmt.Errorf("no metadata found in uploaded object")
	}

	storedMD5, md5Exists := headResult.Metadata["local-md5"]
	storedSHA256, sha256Exists := headResult.Metadata["local-sha256"]

	if !md5Exists || !sha256Exists {
		return fmt.Errorf("integrity metadata not found in uploaded object")
	}

	if storedMD5 != localMD5 || storedSHA256 != localSHA256 {
		return fmt.Errorf("metadata integrity mismatch: stored MD5=%s, local MD5=%s; stored SHA256=%s, local SHA256=%s",
			storedMD5, localMD5, storedSHA256, localSHA256)
	}

	return nil
}
