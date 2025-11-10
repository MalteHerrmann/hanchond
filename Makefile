.phony: format install docs-dev docs-build generate generate-explorer install-deps lint release-dry release

format:
	@gofumpt -l -w .

install:
	@go install

docs-dev:
	@bun i && bun run docs:dev

docs-build:
	@bun i && bun run docs:build

generate:
	@sqlc generate

generate-explorer:
	@sqlc generate -f ./playground/explorer/database/sqlc.yaml

install-deps:
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


lint: lint-go

GOLANGCI_LINT_IMAGE=golangci/golangci-lint:v2.1.6
lint-go:
	@echo "Running golangci-lint..." && \
	docker run -t --rm -v $(CURDIR):/app -w /app $(GOLANGCI_LINT_IMAGE) golangci-lint run

release-dry:
	@goreleaser release --snapshot --clean

release:
	@goreleaser release --skip-validate --clean
