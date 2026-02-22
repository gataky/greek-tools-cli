# Go parameters
BINARY_NAME=greekmaster
MAIN_PATH=cmd/greekmaster/main.go

# Build parameters
LDFLAGS=-ldflags "-s -w"

.PHONY: all build install test clean run tidy help

all: build

## build: Compile the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

## install: Install the binary to $GOPATH/bin
install:
	go install $(LDFLAGS) ./...

## test: Run all tests
test:
	go test -v ./...

## clean: Remove build artifacts and temporary files
clean:
	rm -f $(BINARY_NAME)
	go clean

## run: Build and run the application
run: build
	./$(BINARY_NAME)

## tidy: Clean up go.mod and go.sum
tidy:
	go mod tidy

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
