.PHONY: all build run clean deps install

# Binary name
BINARY=neko

# Build the application
all: build

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download

# Build the binary
build: deps
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) ./cmd/oneko

# Run the application
run: build
	./$(BINARY)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	go clean

# Install system dependencies (Linux)
install-deps-linux:
	@echo "Installing system dependencies for Linux..."
	@if command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update && sudo apt-get install -y libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev; \
	elif command -v dnf >/dev/null 2>&1; then \
		sudo dnf install -y libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel; \
	else \
		echo "Unsupported package manager. Please install X11 development libraries manually."; \
	fi

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test -v ./...

# Check for common issues
vet:
	go vet ./...
