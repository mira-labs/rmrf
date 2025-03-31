SHELL := /bin/bash
GO_BIN := $(shell go env GOPATH)/bin

# Define the output binary name
BINARY := rmrf

.PHONY: all build lint doc install-tools clean

all: build lint doc

# Build the Go project
build:
	go build -o $(BINARY) main.go

# Run Go linters
lint:
	golangci-lint run

# Generate documentation using godox (since godoc -html is deprecated)
doc:
	mkdir -p docs
	$(GO_BIN)/godox -output docs/index.html

# Install required tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/ouqiang/gocron/docs/godox@latest

# Clean up generated files
clean:
	rm -rf $(BINARY) docs
