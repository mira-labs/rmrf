SHELL := /bin/bash

BINARY := rmrf
VERSION := 0.5.0

.PHONY: all build install clean

all: build

build:
	go build -o $(BINARY) ./cmd/rmrf

install: build
	sudo install -m 755 $(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test -v ./...

lint:
	golangci-lint run

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build    - Compile the binary"
	@echo "  install  - Install system-wide"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  clean    - Remove build artifacts"
