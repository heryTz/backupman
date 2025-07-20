## Agent Instructions for `backupman`

### Build, Lint, and Test

- **Go Version:** `1.24`
- **Install/Build:** `go get .`
- **Test:** `go test ./...`
- **Run a single test:** `go test -run ^TestMyFunction$`
- **Lint:** No specific linter is configured. Adhere to `gofmt` standards.

### Code Style Guidelines

- **Formatting:** Use `gofmt` for all Go code.
- **Imports:** Group imports into standard library and third-party blocks.
- **Naming:**
    - `camelCase` for unexported variables and functions.
    - `PascalCase` for exported types, functions, and struct fields.
    - `UPPER_SNAKE_CASE` for constants.
- **Error Handling:**
    - In `main` or initialization, use `log.Fatal(err)`.
    - In other functions, return errors up the call stack.
- **Types:** Use structs for configuration and application state. Use interfaces for abstractions (e.g., `drive.Drive`, `dumper.Dumper`).
- **Dependencies:** Use `go mod tidy` to manage dependencies.
- **Logging:** Use the standard `log` package for logging.
