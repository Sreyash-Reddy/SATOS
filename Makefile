.PHONY: build test clean install run-director run-tl

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=github.com/Sreyash-Reddy/SATOS

# Build targets
BUILD_DIR=bin
DIRECTOR_BIN=$(BUILD_DIR)/satos-director
TL_BIN=$(BUILD_DIR)/satos-tl

# Build both binaries
build: build-director build-tl

build-director:
	$(GOBUILD) -o $(DIRECTOR_BIN) ./cmd/director

build-tl:
	$(GOBUILD) -o $(TL_BIN) ./cmd/tl

# Run tests
test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install binaries to PATH
install: build
	cp $(DIRECTOR_BIN) ~/bin/satos-director
	cp $(TL_BIN) ~/bin/satos-tl

# Run director
run-director:
	$(GOCMD) run ./cmd/director

# Run TL
run-tl:
	$(GOCMD) run ./cmd/tl

# Development helpers
dev-director: build-director
	./$(DIRECTOR_BIN) init

dev-tl: build-tl
	./$(TL_BIN) dispatch

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code
lint:
	golangci-lint run

# Dependency management
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Show help
help:
	@echo "SATOS Build Commands:"
	@echo "  make build         - Build both binaries"
	@echo "  make test          - Run all tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make install       - Install binaries to ~/bin"
	@echo "  make run-director  - Run director"
	@echo "  make run-tl        - Run TL dispatcher"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make help          - Show this help"