APP_NAME := git-community-standards
MAIN_PACKAGE := .
DIST_DIR := dist
VERSION ?= $(shell git describe --tags --always --dirty 2>nul || echo dev)
CURRENT_GOOS := $(shell go env GOOS)
BIN_EXT := $(if $(filter windows,$(CURRENT_GOOS)),.exe,)

LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: help test build release clean

help:
	@echo "Available targets:"
	@echo "  make test     - run Go tests"
	@echo "  make build    - build local binary into ./bin"
	@echo "  make release  - build multi-platform binaries into ./dist"
	@echo "  make clean    - remove build artifacts"

test:
	go test ./...

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)$(BIN_EXT) $(MAIN_PACKAGE)

release:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-arm64.exe $(MAIN_PACKAGE)

clean:
	@rm -rf bin $(DIST_DIR)
