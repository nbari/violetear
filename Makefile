.PHONY: all test clean build cover

GO ?= go

all: build test

build:
	@go build

test: build
	@go test -v

cover:
	@go test -cover && \
	go test -coverprofile=coverage.out  && \
	go tool cover -html=coverage.out
