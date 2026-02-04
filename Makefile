.PHONY: all test clean vendor help

all: test

test:
	@echo "Running framework tests..."
	@go test -v ./...

vendor:
	@echo "Vendoring dependencies..."
	@go mod tidy
	@go mod vendor

clean:
	@echo "Cleaning..."
	@rm -rf vendor/
	@go clean -modcache

help:
	@echo "Ikaitla Framework"
	@echo ""
	@echo "Targets:"
	@echo "  make test   - Run tests"
	@echo "  make vendor - Vendor dependencies"
	@echo "  make clean  - Clean build artifacts"
	@echo "  make help   - Show this help"
