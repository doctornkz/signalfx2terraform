.PHONY: all help vet build install directories

PROJECT_NAME := signalfx2terraform
VENDORFLAGS ?= -mod=vendor

GOFLAGS  = -v

GOPATH   = $(shell go env GOPATH)
GOJUNIT  = $(GOPATH)/bin/go-junit-report
GOCILINT = $(GOPATH)/bin/golangci-lint

MKDIR_P = mkdir -p
INSTALL = install

LDFLAGS += -X "main.version=$(shell date -u '+0.%Y%m%d.%H%M%S')"
GOARCH = amd64

# Detect operating system
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	GOOS=linux
endif
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin
endif

SRC := $(shell find . -type f -name '*.go' -print)

#### DO NOT CHANGE THOSE FOLDERS.
## Folder with all reports [tests,coverate,lint]
FOLDER_REPORT = $(CURDIR)/reports
## Folder with binary [compiled file]
FOLDER_BIN = $(CURDIR)/bin

### COMMANDS
all: build

ship: clean mod build-static install

## Display this help screen
help:
	@awk 'BEGIN {FS = ":.*?##"; printf "Usage: make <target>\n"} /^[a-zA-Z_-]+:.*?##/ {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Directories
directories: $(FOLDER_REPORT)

$(FOLDER_REPORT):
	$(MKDIR_P) $@

## Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
vet:
	go vet ./...

clean: ## Remove previous build
	rm -rf "$(FOLDER_BIN)"
	rm -rf "$(FOLDER_REPORT)"
	@find . -name ".DS_Store" -print0 | xargs -0 rm -f
	go clean -i ./...

build: $(FOLDER_BIN)/$(PROJECT_NAME) ## Build locally binary, will detect your OS and setup the flags

$(FOLDER_BIN)/$(PROJECT_NAME): $(SRC)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -i -o $(FOLDER_BIN)/$(PROJECT_NAME) $(GOFLAGS) -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" ./src

build_macos: $(SRC) ## Build macos binary
	GOOS=darwin GOARCH=amd64 go build -i -o $(FOLDER_BIN)/$(PROJECT_NAME)_macos $(GOFLAGS) -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" ./src

build_linux: $(SRC) ## Build linux binary
	GOOS=linux GOARCH=amd64 go build -i -o $(FOLDER_BIN)/$(PROJECT_NAME)_linux $(GOFLAGS) -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" ./src

build_alpine: $(SRC) ## Build alpine binary
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -i -o $(FOLDER_BIN)/$(PROJECT_NAME)_alpine $(GOFLAGS) -a -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" ./src

install: ## Copy binary to bin folder
	$(INSTALL) -Dm755 $(FOLDER_BIN)/$(PROJECT_NAME) $(GOPATH)/bin/

## Get the dependencies
$(GOCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(LINT_LATEST_VERSION)

$(GOJUNIT):
	(cd /; GO111MODULE=on go get -u github.com/jstemmer/go-junit-report)

mod: $(GOCILINT) $(GOJUNIT)
	go mod vendor

config-git-hooks: ## Config local git-hooks
	git config core.hooksPath .githooks
