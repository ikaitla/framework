.PHONY: all generate build install clean test help

# Default target
all: build

# Generate all code
generate:
	@echo "Generating profile wiring..."
	@go run framework/autodiscovery.go

# Build binary and create symlinks
build: generate
	@echo "Building ikaitla..."
	@mkdir -p bin
	@go build -o bin/ikaitla .
	@echo "Creating profile symlinks..."
	@cd bin && ln -sf ikaitla webkernel 2>/dev/null || true
	@cd bin && ln -sf ikaitla nadim 2>/dev/null || true
	@cd bin && ln -sf ikaitla asnous 2>/dev/null || true
	@cd bin && ln -sf ikaitla wk 2>/dev/null || true
	@cd bin && ln -sf ikaitla nd 2>/dev/null || true
	@cd bin && ln -sf ikaitla as 2>/dev/null || true
	@echo "Build complete: ./bin/"

# Install to user bin directory
install: build
	@echo "Installing to ~/.local/bin/..."
	@mkdir -p ~/.local/bin
	@cp bin/* ~/.local/bin/
	@echo "Installed successfully"

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f generated_profiles.go
	@find cmd -name generated.go -delete
	@echo "Clean complete"

# Show help
help:
	@echo "Ikaitla Build System"
	@echo ""
	@echo "Targets:"
	@echo "  make          - Build everything (default)"
	@echo "  make generate - Generate wiring code"
	@echo "  make build    - Build binary and symlinks"
	@echo "  make install  - Install to ~/.local/bin/"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make help     - Show this help"
