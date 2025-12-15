package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "."
}

// CreateTestSQLiteDB creates a test SQLite database with sample data
func CreateTestSQLiteDB(t *testing.T, dbPath string) {
	// Remove existing database
	os.Remove(dbPath)

	// Create database with test data
	db, err := os.Create(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	db.Close()

	// Use sqlite3 command to create tables and insert data
	cmd := exec.Command("sqlite3", dbPath, `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price DECIMAL(10,2),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		INSERT INTO users (name, email) VALUES 
			('John Doe', 'john@example.com'),
			('Jane Smith', 'jane@example.com'),
			('Bob Johnson', 'bob@example.com');
		
		INSERT INTO products (name, price) VALUES
			('Laptop', 999.99),
			('Mouse', 29.99),
			('Keyboard', 79.99);
	`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create test database schema: %v, output: %s", err, string(output))
	}
}

// CreateTestFile creates a test file with the given content
func CreateTestFile(t *testing.T, dir, filename, content string) string {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	filePath := filepath.Join(dir, filename)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filePath, err)
	}

	return filePath
}

// LogTestInfo logs test information for debugging
func LogTestInfo(t *testing.T, message string, args ...interface{}) {
	t.Logf("[TEST] %s", fmt.Sprintf(message, args...))
}

// AssertFileExists asserts that a file exists
func AssertFileExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Expected file %s to exist, but got error: %v", filePath, err)
	}
}

// AssertFileNotExists asserts that a file does not exist
func AssertFileNotExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		t.Fatalf("Expected file %s to not exist, but it does", filePath)
	}
	if !os.IsNotExist(err) {
		t.Fatalf("Expected file %s to not exist, but got unexpected error: %v", filePath, err)
	}
}

// GetFileContent returns the content of a file
func GetFileContent(t *testing.T, filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}
	return string(content)
}

// CompareFileContent compares the content of two files
func CompareFileContent(t *testing.T, file1, file2 string) {
	content1 := GetFileContent(t, file1)
	content2 := GetFileContent(t, file2)

	if content1 != content2 {
		t.Fatalf("File contents differ:\n--- %s ---\n%s\n\n--- %s ---\n%s",
			file1, content1, file2, content2)
	}
}
