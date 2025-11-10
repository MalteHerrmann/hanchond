.phony: build format install docs-dev docs-build generate generate-explorer lint

export VERSION := $(shell echo $(shell git describe --tags --always --match "v*") | sed 's/^v//')

ld_flags = -X github.com/hanchon/hanchond/cmd.Version=$(VERSION)

build:
	@nix develop -c go build -o build/ -ldflags '$(ld_flags)'

format:
	@nix develop -c golangci-lint fmt -c .golangci.yml

install:
	@nix develop -c go install -ldflags '$(ld_flags)'

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
