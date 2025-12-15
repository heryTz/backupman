//go:build e2e

package tests

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/herytz/backupman/core/drive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MinIO test constants
var (
	minioEndpoint  = "http://localhost:9000"
	minioAccessKey = "minioadmin"
	minioSecretKey = "minioadmin123"
	testBucket     = "test-backups"
	testRegion     = "us-east-1"
)

func TestS3DriveIntegrityE2E(t *testing.T) {
	// Setup MinIO client
	client := setupMinIOClient(t)
	createTestBucket(t, client)
	defer cleanupTestBucket(t, client)

	// Create S3 drive instance
	s3Drive := drive.NewS3Drive(
		"Test S3 Drive",
		testBucket,
		testRegion,
		minioAccessKey,
		minioSecretKey,
		minioEndpoint,
		"integrity-test",
		true,
	)

	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "integrity_test.sql")
	testContent := "This is a test file for S3 integrity verification with content: backup_data_12345"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate expected checksums
	expectedMD5 := calculateMD5Hash(testContent)
	expectedSHA256 := calculateSHA256Hash(testContent)

	// Upload file with integrity verification enabled
	driveFile, err := s3Drive.Upload(testFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, driveFile.Checksum, "Expected checksum to be set in DriveFile")

	// Verify ETag matches MD5 (for single-part uploads)
	etagClean := strings.Trim(driveFile.Checksum, "\"")
	if !strings.Contains(etagClean, "-") {
		assert.Equal(t, expectedMD5, etagClean, "ETag should match MD5 for single-part uploads")
	}

	// Verify metadata in S3
	headResult, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &testBucket,
		Key:    &driveFile.Path,
	})
	assert.NoError(t, err)

	// Check metadata contains our checksums
	require.NotNil(t, headResult.Metadata, "Metadata should not be nil")

	storedMD5 := headResult.Metadata["local-md5"]
	storedSHA256 := headResult.Metadata["local-sha256"]

	assert.Equal(t, expectedMD5, storedMD5, "MD5 metadata should match calculated MD5")
	assert.Equal(t, expectedSHA256, storedSHA256, "SHA256 metadata should match calculated SHA256")
}

func TestS3DriveIntegrityDisabledE2E(t *testing.T) {
	// Setup MinIO client
	client := setupMinIOClient(t)
	createTestBucket(t, client)
	defer cleanupTestBucket(t, client)

	// Create S3 drive with integrity check disabled
	s3Drive := drive.NewS3Drive(
		"Test S3 No Integrity",
		testBucket,
		testRegion,
		minioAccessKey,
		minioSecretKey,
		minioEndpoint,
		"integrity-test-no-check",
		true,
	)
	s3Drive.EnableIntegrityCheck = false

	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "no_integrity_test.sql")
	testContent := "SELECT * FROM no_check_table;"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Upload file without integrity verification
	driveFile, err := s3Drive.Upload(testFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, driveFile.Checksum, "Expected ETag to be set as checksum even with integrity check disabled")

	// Verify no integrity metadata is stored
	headResult, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &testBucket,
		Key:    &driveFile.Path,
	})
	assert.NoError(t, err)

	// Should not have integrity metadata when check is disabled
	if headResult.Metadata != nil {
		assert.Empty(t, headResult.Metadata["local-md5"], "Expected no MD5 metadata when integrity check is disabled")
		assert.Empty(t, headResult.Metadata["local-sha256"], "Expected no SHA256 metadata when integrity check is disabled")
	}
}

// Helper functions
func setupMinIOClient(t *testing.T) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(testRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(minioAccessKey, minioSecretKey, "")),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &minioEndpoint
		o.UsePathStyle = true
	})

	return client
}

func createTestBucket(t *testing.T, client *s3.Client) {
	_, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: &testBucket,
	})
	if err != nil {
		// Bucket might already exist, which is fine for testing
		t.Logf("Bucket creation result: %v", err)
	}
}

func cleanupTestBucket(t *testing.T, client *s3.Client) {
	// List all objects in the bucket
	listOutput, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &testBucket,
	})
	if err != nil {
		t.Logf("Failed to list objects for cleanup: %v", err)
		return
	}

	// Delete all objects
	for _, obj := range listOutput.Contents {
		_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: &testBucket,
			Key:    obj.Key,
		})
		if err != nil {
			t.Logf("Failed to delete object %s: %v", *obj.Key, err)
		}
	}

	// Delete the bucket
	_, err = client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: &testBucket,
	})
	if err != nil {
		t.Logf("Failed to delete bucket: %v", err)
	}
}

func calculateMD5Hash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

func calculateSHA256Hash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
