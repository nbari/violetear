.PHONY: all deps test build cover

GO ?= go

all: build test

deps:
	${GO} get golang.org/x/net/context
	${GO} get github.com/nbari/violetear/middleware

build: deps
build:
	${GO} build

test: deps
test:
	${GO} test

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
