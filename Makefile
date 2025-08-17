.phony: format install docs-dev docs-build generate generate-explorer lint

format:
	@nix develop -c golangci-lint fmt -c .golangci.yml

install:
	@nix develop -c go install

docs-dev:
	@bun i && bun run docs:dev

docs-build:
	@bun i && bun run docs:build

generate:
	@nix develop -c sqlc generate

generate-explorer:
	@nix develop -c sqlc generate -f ./playground/explorer/database/sqlc.yaml

lint:
	@nix develop -c golangci-lint run -c .golangci.yml
