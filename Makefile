.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test

build:
	@go build $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

install:
	@go get $(GOFLAGS) ./...

test: install
	@go test $(GOFLAGS) ./...

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...
