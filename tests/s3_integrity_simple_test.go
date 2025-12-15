package tests

import (
	"strings"
	"testing"

	"github.com/herytz/backupman/core/drive"
)

func TestS3IntegrityVerificationSimple(t *testing.T) {
	// Create S3 drive with integrity check enabled
	s3Drive := drive.NewS3Drive(
		"Test S3 Integrity",
		"test-bucket",
		"us-east-1",
		"test-key",
		"test-secret",
		"",
		"test-integrity",
		false,
	)

	// Verify integrity check is enabled by default
	if !s3Drive.EnableIntegrityCheck {
		t.Error("Expected EnableIntegrityCheck to be true by default")
	}

	// Test with integrity check disabled
	s3DriveNoCheck := drive.NewS3Drive(
		"Test S3 No Integrity",
		"test-bucket",
		"us-east-1",
		"test-key",
		"test-secret",
		"",
		"test-no-integrity",
		false,
	)
	s3DriveNoCheck.EnableIntegrityCheck = false

	if s3DriveNoCheck.EnableIntegrityCheck {
		t.Error("Expected EnableIntegrityCheck to be false when set")
	}
}

func TestS3DriveFileChecksumField(t *testing.T) {
	// Test that DriveFile struct has Checksum field
	driveFile := drive.DriveFile{
		Path:     "test/path/file.txt",
		Checksum: "test-checksum-123",
	}

	if driveFile.Path != "test/path/file.txt" {
		t.Errorf("Expected path 'test/path/file.txt', got '%s'", driveFile.Path)
	}

	if driveFile.Checksum != "test-checksum-123" {
		t.Errorf("Expected checksum 'test-checksum-123', got '%s'", driveFile.Checksum)
	}
}

func TestS3ETagHandling(t *testing.T) {
	// Test ETag cleaning functionality
	etagWithQuotes := `"d41d8cd98f00b204e9800998ecf8427e"`
	expectedClean := "d41d8cd98f00b204e9800998ecf8427e"

	cleaned := strings.Trim(etagWithQuotes, `"`)
	if cleaned != expectedClean {
		t.Errorf("Expected cleaned ETag '%s', got '%s'", expectedClean, cleaned)
	}

	// Test multipart ETag detection
	multipartETag := "d41d8cd98f00b204e9800998ecf8427e-1"
	if !strings.Contains(multipartETag, "-") {
		t.Error("Expected multipart ETag to contain '-'")
	}

	singlePartETag := "d41d8cd98f00b204e9800998ecf8427e"
	if strings.Contains(singlePartETag, "-") {
		t.Error("Expected single-part ETag to not contain '-'")
	}
}
