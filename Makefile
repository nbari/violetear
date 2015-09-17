.PHONY: all test clean build

GO ?= go

all: build test

build:
	@go build

clean:
	@go clean

test: build
	@go test -v
