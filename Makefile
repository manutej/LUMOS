.PHONY: help build test bench clean run lint fmt install

# Variables
BINARY_NAME=lumos
VERSION=$(shell git describe --tags --always 2>/dev/null || echo "dev")
BUILD_DIR=./build
CMD_DIR=./cmd/lumos

help:
	@echo "LUMOS Development Commands"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Build and run with test PDF"
	@echo "  make install      - Install binary globally"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-v       - Run tests with verbose output"
	@echo "  make test-race    - Run tests with race detector"
	@echo "  make coverage     - Generate coverage report"
	@echo ""
	@echo "Quality:"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo ""
	@echo "Performance:"
	@echo "  make bench        - Run benchmarks"
	@echo "  make profile-cpu  - Generate CPU profile"
	@echo "  make profile-mem  - Generate memory profile"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make clean-all    - Remove all artifacts including vendor/"

# Build targets
build:
	@echo "Building $(BINARY_NAME)..."
	go build -v -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: build
	@echo "Building for multiple platforms..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	ls -lh $(BUILD_DIR)/$(BINARY_NAME)-*

install: build
	@echo "Installing $(BINARY_NAME) to ~/bin/..."
	mkdir -p ~/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) ~/bin/$(BINARY_NAME)
	chmod +x ~/bin/$(BINARY_NAME)
	@echo "Installed: ~/bin/$(BINARY_NAME)"
	@echo "Make sure ~/bin is in your PATH"

run: build
	@echo "Running LUMOS..."
	$(BUILD_DIR)/$(BINARY_NAME) --keys

# Testing targets
test:
	@echo "Running tests..."
	go test -v ./...

test-v:
	@echo "Running tests (verbose)..."
	go test -v -count=1 ./...

test-race:
	@echo "Running tests with race detector..."
	go test -race ./...

coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Benchmarking
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./pkg/...

profile-cpu:
	@echo "Profiling CPU usage..."
	go test -cpuprofile=cpu.prof -bench=. ./pkg/...
	go tool pprof -http=:8080 cpu.prof

profile-mem:
	@echo "Profiling memory usage..."
	go test -memprofile=mem.prof -bench=. ./pkg/...
	go tool pprof -http=:8080 mem.prof

# Code quality
fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

# Dependency management
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Code generation (if needed)
generate:
	@echo "Running code generation..."
	go generate ./...

# Cleanup
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f *.prof
	rm -f coverage.out coverage.html

clean-all: clean
	@echo "Cleaning all artifacts..."
	go clean -cache -testcache -modcache
	rm -rf vendor/

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installed development tools"

watch:
	@echo "Watching for changes..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	air -c .air.toml

# CI/CD checks
ci-check: fmt vet test lint
	@echo "CI checks passed!"

# Docker targets (Phase 2+)
docker-build:
	@echo "Building Docker image..."
	docker build -t lumos:latest .

docker-run:
	@echo "Running LUMOS in Docker..."
	docker run -it -v ~/Documents:/data lumos:latest /data/sample.pdf

# Version info
version:
	@echo "Version: $(VERSION)"
	@echo "Go version: $(shell go version)"

# Default target
.DEFAULT_GOAL := help

# Phony targets help
.PHONY: help build build-all install run test test-v test-race coverage bench \
        profile-cpu profile-mem fmt vet lint deps update-deps generate clean \
        clean-all dev-setup watch ci-check docker-build docker-run version
