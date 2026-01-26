PROJECT_DIR := $(CURDIR)
DOCKER_IMG="sm_client:develop"
BIN := "./bin/system_monitoring_daemon"
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S)
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)

# Сборка бинарника демона
.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/daemon

# Сборка образа докер для клиента
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

.PHONY: version
version: build
	$(BIN) version

.PHONY: test
test:
	go test -race -count=100 -timeout=0 ./internal/...

.PHONY: integration-test
integration-test:
	go test -race -count=1 -tags integration ./tests/integration


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.7.2

	
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint