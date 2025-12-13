package tests

import (
	"os"
	"testing"

	"github.com/herytz/backupman/core/drive"
	"github.com/stretchr/testify/assert"
)

func TestNewS3Drive(t *testing.T) {
	label := "test-s3-drive"
	bucket := "test-bucket"
	region := "us-east-1"
	accessKey := "test-access-key"
	secretKey := "test-secret-key"
	endpoint := "https://s3.amazonaws.com"
	prefix := "backups"
	forcePathStyle := false

	s3Drive := drive.NewS3Drive(label, bucket, region, accessKey, secretKey, endpoint, prefix, forcePathStyle)

	assert.NotNil(t, s3Drive)
	assert.Equal(t, label, s3Drive.GetLabel())
	assert.Equal(t, "s3", s3Drive.GetProvider())
	assert.Equal(t, bucket, s3Drive.Bucket)
	assert.Equal(t, region, s3Drive.Region)
	assert.Equal(t, accessKey, s3Drive.AccessKey)
	assert.Equal(t, secretKey, s3Drive.SecretKey)
	assert.Equal(t, endpoint, s3Drive.Endpoint)
	assert.Equal(t, prefix, s3Drive.Prefix)
	assert.Equal(t, forcePathStyle, s3Drive.ForcePathStyle)
}

func TestS3DriveGetLabel(t *testing.T) {
	label := "test-s3-drive"
	s3Drive := drive.NewS3Drive(label, "bucket", "region", "key", "secret", "", "", false)

	assert.Equal(t, label, s3Drive.GetLabel())
}

func TestS3DriveGetProvider(t *testing.T) {
	s3Drive := drive.NewS3Drive("label", "bucket", "region", "key", "secret", "", "", false)

	assert.Equal(t, "s3", s3Drive.GetProvider())
}

// Note: Full integration tests for S3 would require actual AWS credentials and a test bucket.
// These should be run as integration tests, not unit tests.
func TestS3DriveHealthWithInvalidCredentials(t *testing.T) {
	// This test will fail with invalid credentials, which is expected
	s3Drive := drive.NewS3Drive("test", "invalid-bucket", "us-east-1", "invalid-key", "invalid-secret", "", "", false)

	err := s3Drive.Health()
	// We expect this to fail due to invalid credentials
	assert.Error(t, err)
}

// Test file operations with a temporary local file
func TestS3DriveUploadWithInvalidCredentials(t *testing.T) {
	// Create a temporary test file
	tmpFile := "/tmp/test_s3_upload.txt"
	content := "test content for s3 upload"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpFile)

	s3Drive := drive.NewS3Drive("test", "invalid-bucket", "us-east-1", "invalid-key", "invalid-secret", "", "", false)

	// This should fail due to invalid credentials
	_, err = s3Drive.Upload(tmpFile)
	assert.Error(t, err)
}
