GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR=bin
BINARY_NAME=harmony-tui
BINARY_UNIX=$(BINARY_NAME)-unix

version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}@harmony.one

env := GO111MODULE=on

ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}

all: build

build: 
		mkdir -p $(BINARY_DIR)
		$(env) $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v -ldflags="$(ldflags)"  main.go

clean: 
		$(GOCLEAN)
		rm -f $(BINARY_DIR)/$(BINARY_NAME)
		rm -f $(BINARY_DIR)/$(BINARY_UNIX)

run: build
		./$(BINARY_DIR)/$(BINARY_NAME)

deps:
		$(GOGET) github.com/mum4k/termdash
		$(GOGET) github.com/hpcloud/tail

# Cross compilation
build-linux:
		$(env) GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_UNIX) -v -ldflags="$(ldflags)"
