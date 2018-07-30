GITCOMMIT := $(shell git rev-parse HEAD 2>/dev/null)
VERSION := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
BUILDTIME := $(shell TZ=UTC date "+%Y-%m-%d_%H:%M_UTC")

# CLI_SRC_FILES = $(filter-out *_test.go, $(wildcard cmd/client/*.go))
CLI_SRC_FILES = $(filter-out %_test.go, $(wildcard cmd/client/*.go))
SRV_SRC_FILES = $(filter-out %_test.go, $(wildcard cmd/server/*.go))

.PHONY: build-all build-cli build-srv
build-all: build-cli build-srv
build-cli:
	go build -ldflags "-X main.GitCommit=$(GITCOMMIT) -X main.BuildTime=$(BUILDTIME) -X main.Version=$(VERSION)" \
		-v -o bin/gocha-cli ./cmd/client/
build-srv:
	go build -ldflags "-X main.GitCommit=$(GITCOMMIT) -X main.BuildTime=$(BUILDTIME) -X main.Version=$(VERSION)" \
		-v -o bin/gocha-srv ./cmd/server/

.PHONY: tests
tests:
	go test -race ./...

.PHONY: run-cli run-srv
run-cli:
	go run $(CLI_SRC_FILES) || true

run-srv:
	go run $(SRV_SRC_FILES) || true