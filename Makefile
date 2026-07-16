.PHONY: build test test-coverage fmt vet tidy deps lint modern modern-check ci clean help

default: build

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

build: ## Compile the package
	@go build ./...

test: ## Run tests
	@go test ./...

test-coverage: ## Run tests with an HTML coverage report
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

fmt: ## Format code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

tidy: ## Tidy go.mod/go.sum
	@go mod tidy

deps: ## Download dependencies
	@go mod download

lint: ## Run golangci-lint
	@golangci-lint run

modern-check: ## Check for modern Go idioms
	@echo "make: Checking for modern Go code..."
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...

modern: ## Fix modern Go idioms
	@echo "make: Fixing modern Go code..."
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

ci: tidy build vet test ## Run all CI checks locally

clean: ## Remove build/coverage artifacts
	@rm -f coverage.out coverage.html
