#!/bin/bash

# End-to-end test runner for S3 drive with MinIO
set -e

echo "=== Backupman S3 E2E Test Runner ==="

# Check if docker and docker-compose are available
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "Error: docker-compose is not installed or not in PATH"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo "Error: Docker daemon is not running. Please start Docker first."
    exit 1
fi

# Get project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"

# Cleanup function to ensure containers are stopped and removed
cleanup() {
    echo "Cleaning up test containers..."
    cd "$PROJECT_ROOT"
    docker-compose -f docker-compose.test.yml down --volumes --remove-orphans 2>/dev/null || true
    # Only prune images from our test project, not entire system
    docker image prune -f --filter=label=com.docker.compose.project=backupman 2>/dev/null || true
    # Clean up test database file
    rm -f tests/data/test.db 2>/dev/null || true
}

# Set trap for cleanup on script exit
trap cleanup EXIT

# Stop any existing containers
echo "Stopping any existing test containers..."
docker-compose -f docker-compose.test.yml down --volumes --remove-orphans 2>/dev/null || true

# Start MinIO
echo "Starting MinIO..."
docker-compose -f docker-compose.test.yml up -d minio

# Wait for MinIO to be ready
echo "Waiting for MinIO to be ready..."
for i in {1..30}; do
    if curl -f http://localhost:9000/minio/health/live &>/dev/null; then
        echo "MinIO is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "MinIO did not become ready in time"
        exit 1
    fi
    echo "Attempt $i/30: MinIO not ready yet..."
    sleep 2
done

# Create test bucket
echo "Creating test bucket..."
docker-compose -f docker-compose.test.yml exec -T minio sh -c "
    mc alias set minio http://localhost:9000 \$MINIO_ROOT_USER \$MINIO_ROOT_PASSWORD && \
    mc mb minio/test-backups 2>/dev/null || echo 'Bucket may already exist'
" || echo "MinIO bucket setup completed"

# Create test database
echo "Creating test SQLite database..."
mkdir -p tests/data
rm -f tests/data/test.db
sqlite3 tests/data/test.db <<EOF
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
EOF

echo "Test database created at tests/data/test.db"

# Run unit tests
echo "Running unit tests..."
go test -v ./tests/s3_drive_test.go

# Run end-to-end tests
echo "Running end-to-end tests..."
go test -v -tags=e2e ./tests -run "TestS3DriveIntegrityE2E|TestS3DriveIntegrityDisabledE2E"

echo "=== All tests completed successfully! ==="
echo "Containers will be automatically cleaned up..."
