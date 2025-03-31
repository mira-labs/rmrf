# Makefile for Concurrent rm -rf Tool

SHELL := /bin/bash
GO_BIN := $(shell go env GOPATH)/bin
INSTALL_DIR := /usr/local/bin

# Build Configuration
BINARY := rmrf
VERSION := 0.4.1
GOARCH := amd64
MAX_THREADS ?= $(shell getconf _NPROCESSORS_ONLN || echo 4)

# Cross-compilation Targets
PLATFORMS := linux-amd64 windows-amd64 darwin-amd64 darwin-arm64

# Build Flags
LDFLAGS := -X main.version=$(VERSION) -w -s

.PHONY: all install uninstall build test lint coverage docs clean help

.DEFAULT_GOAL := help

## all: Run full CI pipeline (build, test, lint, docs)
all: build test lint docs

## install: Build and install to system bin directory
install: build
	@echo "Installing $(BINARY) to $(INSTALL_DIR)"
	@sudo install -m 755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installation complete. Version:"
	@$(INSTALL_DIR)/$(BINARY) --version

## uninstall: Remove from system bin directory
uninstall:
	@sudo rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "Uninstalled $(BINARY) from $(INSTALL_DIR)"

## build: Compile production binary
build:
	@echo "Building $(BINARY) v$(VERSION)"
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY) main.go

## test: Run tests with race detection
test:
	@go test -race -v ./...

## lint: Run static analysis
lint: install-tools
	@golangci-lint run

## coverage: Generate test coverage report
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## docs: Generate documentation
docs:
	@mkdir -p docs
	@go doc -all > docs/documentation.txt
	@echo "Package documentation generated: docs/documentation.txt"
	@echo "View documentation by running: go doc -http=:6060"

## install-tools: Install required development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/kisielk/errcheck@latest

## clean: Remove build artifacts
clean:
	@rm -rf $(BINARY) $(BINARY)-* coverage* *.out *.html docs/

## help: Show available targets
help:
	@echo "Concurrent rm -rf Tool (v$(VERSION))"
	@echo "Available targets:"
	@awk '/^## / {sub(/^## /, "", $$0); print}' $(MAKEFILE_LIST) | column -t -s ':'