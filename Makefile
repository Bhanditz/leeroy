# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDTAGS=

.PHONY: clean all fmt vet lint build test install static
.DEFAULT: default

all: clean build static fmt lint test vet

build:
	@echo "+ $@"
	@go build -tags "$(BUILDTAGS) cgo" .

static:
	@echo "+ $@"
	CGO_ENABLED=0 go build -tags "$(BUILDTAGS) static_build" -ldflags "-w -extldflags -static" -o leeroy .

fmt:
	@echo "+ $@"
	@gofmt -s -l .

lint:
	@echo "+ $@"
	@golint ./...

test: fmt lint vet
	@echo "+ $@"
	@go test -v -tags "$(BUILDTAGS) cgo" ./...

vet:
	@echo "+ $@"
	@go vet ./...

clean:
	@echo "+ $@"
	@rm -rf leeroy

install:
	@echo "+ $@"
	@go install -v .
