# Backupman Makefile
# Comprehensive build, test, and deployment commands

# Variables
BINARY_NAME=backupman
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_SHA=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commitSHA=$(COMMIT_SHA) -X main.buildDate=$(BUILD_DATE)"
DOCKER_REGISTRY=herytz
DOCKER_IMAGE=$(DOCKER_REGISTRY)/$(BINARY_NAME)

# Directories
BUILD_DIR=build
DIST_DIR=dist
DOCS_DIR=docs
TESTS_DIR=tests

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golint

# Default target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
.PHONY: deps
deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: fmt
fmt: ## Format Go code
	$(GOFMT) -s -w .

.PHONY: vet
vet: ## Run go vet
	$(GOCMD) vet ./...

.PHONY: lint
lint: ## Run golint (requires golint to be installed)
	@which golint > /dev/null || (echo "golint not installed. Run: go install golang.org/x/lint/golint@latest" && exit 1)
	golint ./...

.PHONY: check
check: fmt vet lint ## Run all code quality checks

# Build commands
.PHONY: build
build: ## Build the binary for current platform
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: build-all
build-all: ## Build binaries for all platforms
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

.PHONY: build-release
build-release: ## Build release binaries using GoReleaser
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest" && exit 1)
	goreleaser build --snapshot --clean

.PHONY: release
release: ## Create a full release using GoReleaser
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest" && exit 1)
	goreleaser release --clean

# Test commands
.PHONY: test
test: ## Run unit tests
	$(GOTEST) -v ./$(TESTS_DIR)

.PHONY: test-short
test-short: ## Run short tests only
	$(GOTEST) -v -short ./$(TESTS_DIR)

.PHONY: test-integration
test-integration: ## Run integration tests
	$(GOTEST) -v -tags=integration ./$(TESTS_DIR)

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests (requires MinIO setup)
	@echo "Running E2E tests..."
	@bash $(TESTS_DIR)/run_s3_e2e_tests.sh

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@mkdir -p $(BUILD_DIR)
	$(GOTEST) -v -coverprofile=$(BUILD_DIR)/coverage.out ./$(TESTS_DIR)
	$(GOCMD) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report generated: $(BUILD_DIR)/coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detector
	$(GOTEST) -v -race ./$(TESTS_DIR)

.PHONY: test-all
test-all: test test-integration test-e2e-quick ## Run all tests

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):latest .
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):$(VERSION)

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	docker push $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(VERSION)

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):latest

.PHONY: docker-compose-up
docker-compose-up: ## Start development services with Docker Compose
	docker-compose -f compose.yml up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop development services
	docker-compose -f compose.yml down

.PHONY: docker-compose-logs
docker-compose-logs: ## Show Docker Compose logs
	docker-compose -f compose.yml logs -f

# Documentation commands
.PHONY: docs-install
docs-install: ## Install documentation dependencies
	cd $(DOCS_DIR) && pnpm install

.PHONY: docs-dev
docs-dev: ## Start documentation development server
	cd $(DOCS_DIR) && pnpm start

.PHONY: docs-build
docs-build: ## Build documentation
	cd $(DOCS_DIR) && pnpm build

.PHONY: docs-serve
docs-serve: ## Serve built documentation
	cd $(DOCS_DIR) && pnpm serve

.PHONY: docs-deploy
docs-deploy: ## Deploy documentation
	cd $(DOCS_DIR) && pnpm deploy

.PHONY: docs-clear
docs-clear: ## Clear documentation cache
	cd $(DOCS_DIR) && pnpm clear

# Database commands
.PHONY: db-setup
db-setup: ## Set up test databases
	docker-compose -f compose.yml up -d postgres mariadb
	@echo "Waiting for databases to be ready..."
	@sleep 10
	@echo "Databases are ready!"

.PHONY: db-reset
db-reset: ## Reset test databases
	docker-compose -f compose.yml down -v
	docker-compose -f compose.yml up -d postgres mariadb

.PHONY: db-migrate
db-migrate: ## Run database migrations
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	./$(BUILD_DIR)/$(BINARY_NAME) migrate

# Development environment
.PHONY: dev-setup
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@echo "Installing Go dependencies..."
	$(MAKE) deps
	@echo "Setting up documentation..."
	$(MAKE) docs-install
	@echo "Starting development services..."
	$(MAKE) docker-compose-up
	@echo "Development environment ready!"

.PHONY: dev-run
dev-run: ## Run the application in development mode
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	./$(BUILD_DIR)/$(BINARY_NAME) serve

.PHONY: dev-health
dev-health: ## Check application health
	curl -f http://localhost:8080/health || echo "Application is not running or unhealthy"

# Cleanup commands
.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out

.PHONY: clean-docker
clean-docker: ## Clean Docker resources
	docker system prune -f
	docker volume prune -f

.PHONY: clean-all
clean-all: clean clean-docker ## Clean everything

# Utility commands
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT_SHA)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Go Version: $(GO_VERSION)"

.PHONY: install
install: build ## Install binary to /usr/local/bin
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

.PHONY: uninstall
uninstall: ## Uninstall binary from /usr/local/bin
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# CI/CD commands
.PHONY: ci
ci: deps check test-all build ## Run CI pipeline locally

.PHONY: pre-commit
pre-commit: fmt vet test-short ## Run pre-commit checks

# Kubernetes commands
.PHONY: k8s-deploy
k8s-deploy: ## Deploy to Kubernetes (requires kubectl)
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/

.PHONY: k8s-logs
k8s-logs: ## Show Kubernetes logs
	kubectl logs -f deployment/backupman

.PHONY: k8s-status
k8s-status: ## Show Kubernetes deployment status
	kubectl get pods -l app=backupman

# Security commands
.PHONY: security-scan
security-scan: ## Run security scan (requires gitleaks)
	@which gitleaks > /dev/null || (echo "gitleaks not installed. Visit: https://github.com/gitleaks/gitleaks" && exit 1)
	gitleaks detect

.PHONY: security-audit
security-audit: ## Audit Go dependencies
	$(GOMOD) verify
	@which gosec > /dev/null || (echo "gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...

# Performance commands
.PHONY: bench
bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./$(TESTS_DIR)

.PHONY: profile
profile: ## Generate CPU profile
	$(GOTEST) -cpuprofile=$(BUILD_DIR)/cpu.prof -bench=. ./$(TESTS_DIR)
	@echo "CPU profile generated: $(BUILD_DIR)/cpu.prof"

# Monitoring commands
.PHONY: monitor
monitor: ## Start monitoring dashboard (requires monitoring tools)
	@echo "Starting monitoring..."
	@echo "This would typically start Grafana/Prometheus dashboards"

# Backup commands (project-specific)
.PHONY: backup-test
backup-test: ## Run a test backup
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	./$(BUILD_DIR)/$(BINARY_NAME) run-backup --config config-example.yml

.PHONY: backup-list
backup-list: ## List all backups
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	./$(BUILD_DIR)/$(BINARY_NAME) list-backups

.PHONY: backup-retry
backup-retry: ## Retry failed backups
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	./$(BUILD_DIR)/$(BINARY_NAME) retry-backups
