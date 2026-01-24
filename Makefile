PROJECT_DIR := $(CURDIR)
BIN := "./bin/sm"
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S)
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)

.PHONY: build
build:
# 	go build -v -o ./bin/sm_app/daemon -ldflags "$(LDFLAGS)" ./cmd/daemon
# 	go build -v -o ./bin/sm_app/client -ldflags "$(LDFLAGS)" ./cmd/client
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/sm

.PHONY: run
run: build
	LOG_LEVEL=DEBUG $(BIN) -config ./config/config.yaml

.PHONY: version
version: build
	$(BIN) version


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.7.2

	
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint