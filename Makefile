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

lint:
	@golangci-lint run --fix --out-format=line-number --issues-exit-code=0 --config .golangci.yml --color always ./...

release-dry:
	@goreleaser release --snapshot --clean

release:
	@goreleaser release --skip-validate --clean
