.phony: build format install docs-dev docs-build generate generate-explorer lint tools

export VERSION := $(shell echo $(shell git describe --tags --always --match "v*") | sed 's/^v//')

ld_flags = -X github.com/hanchon/hanchond/cmd.Version=$(VERSION)

build:
	@go build -o build/ -ldflags '$(ld_flags)'

format:
	@golangci-lint fmt -c .golangci.yml

install:
	@go install -ldflags '$(ld_flags)'

docs-dev:
	@bun i && bun run docs:dev

docs-build:
	@bun i && bun run docs:build

generate:
	@sqlc generate

generate-explorer:
	@sqlc generate -f ./playground/explorer/database/sqlc.yaml

lint:
	@golangci-lint run -c .golangci.yml

lint-fix:
	@golangci-lint run -c .golangci.yml --fix

tools:
	@mise install
