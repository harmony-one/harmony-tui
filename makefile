# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=harmony-tui
BINARY_UNIX=$(BINARY_NAME)-unix

all: build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
#test: 
#        $(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/mum4k/termdash
		$(GOGET) github.com/hpcloud/tail

# Cross compilation
build-linux:
		GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v