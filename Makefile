DOCKER_IMG="sm_app:dev"

PROJECT_DIR := $(CURDIR)

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o ./bin/sm_app/daemon -ldflags "$(LDFLAGS)" ./cmd/daemon

run: build
	$(BIN) -config ./configs/config.yaml

version: build
	$(BIN) version


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.7.2

	
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint