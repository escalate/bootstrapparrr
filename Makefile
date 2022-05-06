SHELL = /bin/bash
.SHELLFLAGS = -e -o pipefail -c
PROJECT_NAME := "bootstrapparrr"

.PHONY: build
build:
	go build .

.PHONY: run
run:
	go run main.go

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm --force $(PROJECT_NAME)
