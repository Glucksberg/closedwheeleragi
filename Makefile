# Coder AGI Makefile
# Cross-platform build system

.PHONY: all build run test clean install deps

# Binary name
BINARY_NAME := ClosedWheeler

# Detect OS
ifeq ($(OS),Windows_NT)
    BINARY_EXT := .exe
    RM := del /Q
    MKDIR := mkdir
else
    BINARY_EXT :=
    RM := rm -f
    MKDIR := mkdir -p
endif

BINARY := $(BINARY_NAME)$(BINARY_EXT)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -s -w

all: deps build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build the binary
build:
	@echo "🔨 Building $(BINARY)..."
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/agi

# Build for all platforms
build-all: deps
	@echo "🔨 Building for all platforms..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/agi
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/agi
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/agi
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/agi

# Run the application
run: deps
	@echo "🚀 Running Coder AGI..."
	$(GORUN) ./cmd/agi -project workplace

# Run with specific project
run-project:
	@echo "🚀 Running Coder AGI on $(PROJECT)..."
	$(GORUN) ./cmd/agi -project $(PROJECT)

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(RM) $(BINARY)
	$(RM) coverage.out coverage.html

# Install to GOPATH/bin
install: build
	@echo "📥 Installing to GOPATH/bin..."
	$(GOCMD) install ./cmd/agi

# Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run ./...

# Help
help:
	@echo "Coder AGI Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  deps          Install dependencies"
	@echo "  build         Build the binary"
	@echo "  build-all     Build for all platforms (Windows, macOS, Linux)"
	@echo "  run           Run the application"
	@echo "  run-project   Run with specific project (make run-project PROJECT=/path)"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  clean         Clean build artifacts"
	@echo "  install       Install to GOPATH/bin"
	@echo "  fmt           Format code"
	@echo "  lint          Lint code"
	@echo "  help          Show this help"
