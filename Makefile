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
	go build -v -o ./dist/ ./...

.PHONY: run
run:
	go run main.go

.PHONY: lint-dockerfile
lint-dockerfile:
	find $(PWD) -name Dockerfile* -print0 | xargs -0 -I % hadolint %

.PHONY: build-docker-image
build-docker-image:
	docker-compose build

.PHONY: run-docker-image
run-docker-image:
	docker-compose up
