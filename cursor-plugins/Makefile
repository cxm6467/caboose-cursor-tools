.PHONY: build test clean install

# Build all binaries
build:
	@echo "Building cursor-plugins..."
	@mkdir -p bin
	go build -o bin/parse-session ./cmd/parse-session
	@echo "✓ Built parse-session"
	go build -o bin/lint-daemon ./cmd/lint-daemon
	@echo "✓ Built lint-daemon"
	go build -o bin/watch-ci ./cmd/watch-ci
	@echo "✓ Built watch-ci"
	go build -o bin/watch-copilot-reviews ./cmd/watch-copilot-reviews
	@echo "✓ Built watch-copilot-reviews"
	# Add more as we port them
	# go build -o bin/cursor-plugin ./cmd/cursor-plugin

# Run tests
test:
	go test -v ./...

# Run tests with coverage
cover:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install binaries to /usr/local/bin
install: build
	@echo "Installing to /usr/local/bin..."
	sudo cp bin/* /usr/local/bin/
	@echo "✓ Installed"

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...

# Tidy dependencies
tidy:
	go mod tidy

# Development build with race detection
dev:
	go build -race -o bin/parse-session ./cmd/parse-session
