# ============================================================
# Makefile for golang-adk-exploration-1
# Google ADK v2 Go Agent Project
# ============================================================

# ── Project config ──────────────────────────────────────────
MODULE      := github.com/akhapre-dev/golang-exploration-1
CMD         := ./cmd/agent
BINARY_NAME := agent
BINARY_DIR  := bin
BINARY      := $(BINARY_DIR)/$(BINARY_NAME)

# Secrets / environment
ENV_FILE    ?= .env
IMAGE_REPO  ?= $(BINARY_NAME)

# Go tools
GO          := go
GOFLAGS     :=
GOTEST      := $(GO) test
GOBUILD     := $(GO) build
GOVET       := $(GO) vet
GOFMT       := gofmt

# Build metadata injected via ldflags
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Linting
GOLANGCI    := golangci-lint
GOLANGCI_VERSION := v1.64.8

# Test options
TEST_FLAGS  ?= -v -race -count=1
COVER_DIR   := coverage
COVER_OUT   := $(COVER_DIR)/coverage.out
COVER_HTML  := $(COVER_DIR)/coverage.html

# ── Colours ─────────────────────────────────────────────────
CYAN   := \033[36m
RESET  := \033[0m
BOLD   := \033[1m

.DEFAULT_GOAL := help
.PHONY: help build build-all run run-web clean test test-short test-cover cover-html \
        lint lint-fix fmt fmt-check vet tidy deps deps-upgrade deps-verify check \
        install install-tools \
        docker-build docker-run docker-run-web docker-push docker-clean

# ── Help ────────────────────────────────────────────────────
help: ## Show this help message
	@printf '$(BOLD)Usage:$(RESET)\n  make $(CYAN)<target>$(RESET)\n\n'
	@printf '$(BOLD)Targets:$(RESET)\n'
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { \
		printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ── Build ────────────────────────────────────────────────────
build: ## Compile the agent binary to bin/agent
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(GOFLAGS) $(LDFLAGS) -o $(BINARY) $(CMD)
	@echo "✅  Built $(BINARY) (version=$(VERSION) commit=$(COMMIT))"

build-all: ## Cross-compile for linux/amd64 and darwin/arm64
	@mkdir -p $(BINARY_DIR)
	GOOS=linux  GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64   $(CMD)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-arm64  $(CMD)
	@echo "✅  Cross-compiled binaries in $(BINARY_DIR)/"

# ── Run ─────────────────────────────────────────────────────
run: ## Run the agent in interactive CLI mode (loads GOOGLE_API_KEY from .env)
	@test -f $(ENV_FILE) || (echo "⚠️  $(ENV_FILE) not found — copy .env.example to .env and set GOOGLE_API_KEY"; exit 1)
	@export $$(grep -v '^\s*#' $(ENV_FILE) | grep -v '^\s*$$' | xargs) && $(GO) run $(CMD)

run-web: ## Run the agent with the ADK Web UI (loads GOOGLE_API_KEY from .env)
	@test -f $(ENV_FILE) || (echo "⚠️  $(ENV_FILE) not found — copy .env.example to .env and set GOOGLE_API_KEY"; exit 1)
	@export $$(grep -v '^\s*#' $(ENV_FILE) | grep -v '^\s*$$' | xargs) && $(GO) run $(CMD) web webui api

# ── Test ─────────────────────────────────────────────────────
test: ## Run all tests
	$(GOTEST) $(TEST_FLAGS) ./...

test-short: ## Run tests excluding long-running/integration tests
	$(GOTEST) -short -race -count=1 ./...

test-cover: ## Run tests with coverage report
	@mkdir -p $(COVER_DIR)
	$(GOTEST) $(TEST_FLAGS) -coverprofile=$(COVER_OUT) -covermode=atomic ./...
	@$(GO) tool cover -func=$(COVER_OUT) | tail -1

cover-html: test-cover ## Generate HTML coverage report and open in browser
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo "📊  Coverage report: $(COVER_HTML)"
	@open $(COVER_HTML) 2>/dev/null || xdg-open $(COVER_HTML) 2>/dev/null || true

# ── Code quality ─────────────────────────────────────────────
vet: ## Run go vet on all packages
	$(GOVET) ./...

fmt: ## Format all Go source files with gofmt
	$(GOFMT) -w -s .

fmt-check: ## Check formatting without modifying files (CI-friendly)
	@out=$$($(GOFMT) -l -s .); \
	if [ -n "$$out" ]; then \
		echo "❌  The following files are not formatted:"; \
		echo "$$out"; \
		exit 1; \
	fi
	@echo "✅  All files are properly formatted"

lint: ## Run golangci-lint with project config
	$(GOLANGCI) run ./...

lint-fix: ## Run golangci-lint and auto-fix issues where possible
	$(GOLANGCI) run --fix ./...

check: fmt-check vet lint test ## Run all checks: fmt, vet, lint, test (CI gate)

# ── Dependencies ─────────────────────────────────────────────
tidy: ## Tidy go.mod and go.sum
	$(GO) mod tidy

deps: ## Download all module dependencies
	$(GO) mod download

deps-upgrade: ## Upgrade all direct dependencies to their latest versions
	$(GO) get -u ./...
	$(GO) mod tidy

deps-verify: ## Verify module dependencies against go.sum
	$(GO) mod verify

# ── Install ──────────────────────────────────────────────────
install: build ## Install the agent binary to $GOPATH/bin
	$(GO) install $(LDFLAGS) $(CMD)
	@echo "✅  Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

install-tools: ## Install development tools (golangci-lint)
	@echo "📦  Installing golangci-lint $(GOLANGCI_VERSION)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_VERSION)
	@echo "✅  Tools installed"

# ── Docker ───────────────────────────────────────────────────
# Secrets (GOOGLE_API_KEY, etc.) live in .env (gitignored).
# They are passed to the container at runtime only — NEVER baked into the image.

docker-build: ## Build the Docker image (multi-stage; secrets are NOT embedded)
	@test -f Dockerfile || (echo "❌  Dockerfile not found"; exit 1)
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(IMAGE_REPO):$(VERSION) \
		-t $(IMAGE_REPO):latest \
		.
	@echo "✅  Docker image built: $(IMAGE_REPO):$(VERSION)"

docker-run: docker-build ## Build image then run agent CLI (injects GOOGLE_API_KEY from .env)
	@test -f $(ENV_FILE) || (echo "❌  $(ENV_FILE) not found — copy .env.example to .env and set your API key"; exit 1)
	docker run --rm -it \
		--env-file $(ENV_FILE) \
		$(IMAGE_REPO):$(VERSION)

docker-run-web: docker-build ## Build image then run ADK Web UI on port 8080 (injects GOOGLE_API_KEY from .env)
	@test -f $(ENV_FILE) || (echo "❌  $(ENV_FILE) not found — copy .env.example to .env and set your API key"; exit 1)
	docker run --rm -it \
		--env-file $(ENV_FILE) \
		-p 8080:8080 \
		$(IMAGE_REPO):$(VERSION) web

docker-push: ## Push image to a remote registry (set IMAGE_REPO=registry/image before calling)
	docker push $(IMAGE_REPO):$(VERSION)
	docker push $(IMAGE_REPO):latest
	@echo "✅  Pushed $(IMAGE_REPO):$(VERSION)"

docker-clean: ## Remove local Docker images for this project
	docker rmi -f $(IMAGE_REPO):$(VERSION) $(IMAGE_REPO):latest 2>/dev/null || true
	@echo "🧹  Removed Docker images"

# ── Cleanup ──────────────────────────────────────────────────
clean: ## Remove build artifacts and coverage reports
	@rm -rf $(BINARY_DIR) $(COVER_DIR)
	@echo "🧹  Cleaned build artifacts"
