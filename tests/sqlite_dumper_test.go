package tests

import (
	"database/sql"
	"os"
	"path"
	"testing"

	"github.com/herytz/backupman/core/dumper"
	_ "github.com/mattn/go-sqlite3"
)

func TestSqliteDumper(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	sourceDbPath := path.Join(tmpDir, "source.db")
	dumpFolder := path.Join(tmpDir, "dumps")

	// Create a test SQLite database with some data
	db, err := sql.Open("sqlite3", sourceDbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create a test table and insert some data
	_, err = db.Exec(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			value INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO test_table (name, value) VALUES
		('test1', 100),
		('test2', 200),
		('test3', 300)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	db.Close()

	// Create the dumper
	sqliteDumper := dumper.NewSqliteDumper("Test SQLite", dumpFolder, sourceDbPath)

	// Test Health check
	err = sqliteDumper.Health()
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Test GetLabel
	if sqliteDumper.GetLabel() != "Test SQLite" {
		t.Errorf("Expected label 'Test SQLite', got '%s'", sqliteDumper.GetLabel())
	}

	// Test Dump
	dumpPath, err := sqliteDumper.Dump()
	if err != nil {
		t.Fatalf("Dump failed: %v", err)
	}

	// Verify the dump file exists
	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		t.Errorf("Dump file does not exist at %s", dumpPath)
	}

	// Verify the dump file has content
	fileInfo, err := os.Stat(dumpPath)
	if err != nil {
		t.Fatalf("Failed to stat dump file: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Error("Dump file is empty")
	}

	// Open the dump file and verify the data
	dumpDb, err := sql.Open("sqlite3", dumpPath)
	if err != nil {
		t.Fatalf("Failed to open dump database: %v", err)
	}
	defer dumpDb.Close()

	// Check that the table exists
	var tableName string
	err = dumpDb.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&tableName)
	if err != nil {
		t.Fatalf("Test table not found in dump: %v", err)
	}

	// Check that the data was copied
	var count int
	err = dumpDb.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count rows in dump: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 rows in dump, got %d", count)
	}

	// Verify one of the rows
	var name string
	var value int
	err = dumpDb.QueryRow("SELECT name, value FROM test_table WHERE id = 1").Scan(&name, &value)
	if err != nil {
		t.Fatalf("Failed to query row from dump: %v", err)
	}
	if name != "test1" || value != 100 {
		t.Errorf("Expected name='test1' value=100, got name='%s' value=%d", name, value)
	}

	t.Logf("Successfully dumped SQLite database to %s", dumpPath)
}
