COVERAGE_OUT := .coverage.out
COVERAGE_HTML := coverage.html
GO_EXEC := $(shell which go)
GO_TEST_COVER := $(GO_EXEC) test -shuffle=on -cover -covermode=atomic
BIN_NAME := arangom

default: build

.PHONY: help
help: ## Show available targets
	@echo "Available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: dep
dep: ## Download dependencies
	$(GO_EXEC) mod tidy
	$(GO_EXEC) mod download

.PHONY: dep-update
dep-update: ## Update dependencies
	$(GO_EXEC) get -t -u ./...
	$(GO_EXEC) mod tidy
	$(GO_EXEC) mod download

.PHONY: build
build: dep ## Build binary
	goreleaser build --clean --snapshot --single-target
	@find bin -name "$(BIN_NAME)" -exec cp "{}" bin/ \;

.PHONY: release
release: lint test ## Release a new version on GitHub
	goreleaser release --clean

.PHONY: bench
bench: ## Run benchmarks
	$(GO_EXEC) test -bench=. -benchmem -benchtime=10s ./...

.PHONY: format
format: dep ## Format source code
	gofmt -l -s -w $(shell pwd)
	goimports -w $(shell pwd)

.PHONY: lint
lint: ## Run linters
	golangci-lint run --timeout 5m

.PHONY: test
test: ## Run unit tests
	@rm -f $(COVERAGE_OUT)
	$(GO_TEST_COVER) -race -coverprofile=$(COVERAGE_OUT) ./...

.PHONY: coverage.html
coverage.html: ## Generate html coverage report from previous test run
	$(GO_EXEC) tool cover -html "$(COVERAGE_OUT)" -o "$(COVERAGE_HTML)"

.PHONY: changelog
changelog: ## Generate changelog
	git cliff > CHANGELOG.md
