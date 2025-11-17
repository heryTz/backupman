# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Backupman is a minimalist backup system for databases written in Go. It creates database dumps, stores backups locally or in cloud services (Google Drive, local storage), automates scheduled backups with retention rules, and sends notifications about backup status.

**Key capabilities:**
- Multi-database support: MySQL, PostgreSQL (in progress), SQLite
- Hybrid storage: Local filesystem, Google Drive, S3 (planned)
- Cron-based scheduled backups
- Email and webhook notifications
- CLI and HTTP API interfaces
- Time-based retention with auto-cleanup

## Build and Development

**Go Version:** 1.24

**Common Commands:**
```bash
# Install dependencies
go get .

# Run all tests
go test ./...

# Run a specific test
go test -run ^TestBackupRetry$ ./tests

# Build the binary
go build -o backupman .

# Run CLI commands
./backupman run -c ./config.yml          # Run backup once
./backupman serve -p 8080 -c ./config.yml # Start HTTP server
./backupman retry <backup-id> -c ./config.yml
./backupman health -c ./config.yml
./backupman version
./backupman auth-google --client-secret-file ./google-client-secret.json --token-file ./google-token.json
```

**Docker:**
```bash
# Build and push handled by GoReleaser in CI
docker run herytz/backupman:latest
```

**Release:**
- Releases are automated via GitHub Actions using GoReleaser
- Triggered manually via workflow_dispatch
- Version info is injected via ldflags into main.go variables

## Architecture

### Core Structure

The codebase follows a layered architecture with clear separation of concerns:

```
main.go              # Entry point, CLI setup with Cobra
├── cmd/             # Command implementations (run, serve, retry, health, version, auth-google)
├── core/            # Business logic layer
│   ├── application/ # App struct and dependency injection container
│   ├── service/     # Business logic (backup, retry, health, cleanup)
│   ├── dao/         # Data access layer with interfaces
│   ├── dumper/      # Database dump implementations (MySQL, etc.)
│   ├── drive/       # Storage provider implementations (local, Google Drive)
│   ├── notifier/    # Notification implementations (email, webhook)
│   ├── mailer/      # Email transport abstraction
│   ├── model/       # Data models (Backup, DriveFile)
│   └── lib/         # Shared utilities (DB connections, health checks, URLs)
├── http/            # HTTP server (Gin) with API routes and scheduler
├── migration/       # Database migrations (MySQL)
└── tests/           # Integration and unit tests
```

### Key Design Patterns

**1. Dependency Injection via App Container**
- `core/application/app.go` defines the `App` struct that holds all dependencies
- `NewApp(config AppConfig)` constructs the app by wiring together drives, dumpers, DAOs, notifiers
- Commands receive the `App` and pass it to services

**2. Interface-Based Abstractions**
- `dumper.Dumper` interface allows different database dump strategies
- `drive.Drive` interface abstracts storage providers (local, Google Drive, future S3)
- `dao.BackupDao` and `dao.DriveFileDao` enable swappable persistence layers (MySQL, in-memory)
- `notifier.Notifier` interface for multiple notification channels
- `mailer.Mailer` interface for email transport

**3. Database Providers**
- **MySQL:** Default production database using `dao/mysql/*`
- **Memory:** In-memory database for development/testing using `dao/memory/*`
- Configured via `config.yml` (`database.provider: mysql` or `database.provider: memory`)

**4. Dual Execution Mode**
- `APP_MODE_CLI`: Synchronous backup execution, blocking retention cleanup
- `APP_MODE_WEB`: HTTP server mode with async cleanup (goroutines) and cron scheduling

### Backup Flow

1. **Initiation:** CLI `run` command or HTTP POST `/api/backups` triggers `service.Backup(app)`
2. **Database Dump:** For each configured dumper (data source), create a dump file locally
3. **Storage Upload:** For each configured drive, upload the dump file
4. **Tracking:** Backup status and DriveFile records stored in the internal database
5. **Post-Backup:** Execute `AfterBackup()` to finalize backup state and notify
6. **Retention Cleanup:** If enabled, `RemoveOldBackup()` deletes backups older than configured days

### Configuration

`config.yml` is the single source of configuration. Key sections:

- `database`: Internal DB for tracking backups (mysql or memory provider)
- `data_sources`: List of databases to backup (MySQL, PostgreSQL, etc.)
- `drives`: Storage destinations (local, google_drive)
- `notifiers`: Email (SMTP) and webhook configurations
- `retention`: Auto-cleanup based on age (in days)
- `http`: API server settings, API keys, cron schedule for automated backups

**Important:** HTTP server includes the cron automation feature via `http.SetupScheduler()`. This will eventually be decoupled into a standalone daemon.

### Google Drive Authentication

Google Drive requires OAuth2 authentication:
1. Obtain `google-client-secret.json` from Google Cloud Console
2. Run `./backupman auth-google --client-secret-file ./google-client-secret.json --token-file ./google-token.json --open-url`
3. Follow the browser OAuth flow to generate `google-token.json`
4. Configure `drives` in `config.yml` with paths to both files

### HTTP API

Routes (require `Authorization: Bearer <api_key>` header):
- `GET /health` - Health check (no auth required)
- `GET /ping` - Ping endpoint (no auth required)
- `GET /api/backups` - List all backups
- `POST /api/backups` - Trigger a backup job
- `GET /api/backups/:id/generate-download-url` - Generate download URL
- `GET /api/backups/:id/download` - Download backup file

Scheduler: When `http.backup_job.enabled: true`, the server starts a cron job to run backups automatically.

## Code Style

- **Formatting:** Use `gofmt` for all Go code
- **Imports:** Group into standard library and third-party blocks
- **Naming:**
  - `camelCase` for unexported variables/functions
  - `PascalCase` for exported types/functions/fields
  - `UPPER_SNAKE_CASE` for constants (e.g., `BACKUP_STATUS_PENDING`)
- **Error Handling:**
  - In `main` or initialization: `log.Fatal(err)`
  - In other functions: return errors up the call stack
- **Logging:** Use standard `log` package; avoid adding new logging libraries

## Testing

- Tests located in `/tests` directory
- Use `tests/app_mock.go` for creating test fixtures
- Integration tests require MySQL (see `backup_dao_mysql_test.go`)
- Mock implementations available for drives, dumpers, mailers, notifiers
- Run `go test ./...` before committing

## Database Migrations

MySQL migrations are in `migration/mysql/`:
- Written as Go code, not SQL files
- Applied automatically via `migration.RunMigration()` on startup when using MySQL provider
- Create new migration files following the pattern `N_description.go`

## Current Development

**Branch:** `36-postgres-db-support` - Adding PostgreSQL database support
- Implement `dumper.PostgresDumper` (similar to `dumper.MysqlDumper`)
- Update `core/application/app.go` to handle PostgreSQL data source config
- Add PostgreSQL test coverage

## Project Dependencies

Key external libraries:
- `spf13/cobra` - CLI framework
- `gin-gonic/gin` - HTTP framework
- `go-co-op/gocron/v2` - Cron scheduling
- `go-sql-driver/mysql` - MySQL driver
- `golang.org/x/oauth2`, `google.golang.org/api` - Google Drive integration
- `goccy/go-yaml` - YAML config parsing
- `emersion/go-smtp` - Email sending
