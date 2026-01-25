PROJECT_DIR := $(CURDIR)
BIN := "./bin/system_monitoring_daemon"
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S)
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)

.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/daemon

.PHONY: run
run: build
	LOG_LEVEL=DEBUG $(BIN) -config ./config/config.yaml

.PHONY: version
version: build
	$(BIN) version

.PHONY: test
test:
	go test -race -count=100 ./...

.PHONY: test-integr
integration-test:
	go test -race -count=1 -tags integration ./tests/integration


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.7.2

	
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint