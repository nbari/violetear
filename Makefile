.PHONY: all test build cover

GO ?= go

all: build test

build:
	${GO} build

test:
	${GO} test -v

cover:
	${GO} test -cover && \
	go test -coverprofile=coverage.out  && \
	go tool cover -html=coverage.out
