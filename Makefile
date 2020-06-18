VERSION         := 0.1.0
SHORT_COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
GO_VERSION      := $(shell go version | awk '{ print $$3}' | sed 's/^go//')

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

LD_FLAGS_PKG ?= github.com/fenrirunbound/kubeconfig-merge/
LD_FLAGS :=
LD_FLAGS +=  -X "$(LD_FLAGS_PKG).version=$(VERSION)"
LD_FLAGS +=  -X "$(LD_FLAGS_PKG).commit=$(SHORT_COMMIT)"
LD_FLAGS +=  -X "$(LD_FLAGS_PKG).goVersion=$(GO_VERSION)"

.PHONY: build
build:
	go install -ldflags '$(LD_FLAGS)' ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./...