SHELL = /bin/bash
.SHELLFLAGS = -e -o pipefail -c
.DEFAULT_GOAL := build

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: build
build:
	goreleaser build --single-target --snapshot --rm-dist

.PHONY: run
run:
	go run main.go
