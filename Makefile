.PHONY: all deps clean test build cover

GO ?= go

all: build test

deps:
	${GO} get github.com/nbari/violetear/middleware

build: deps
build:
	${GO} build

clean:
	@rm -rf *.out

test: deps
test:
	${GO} test

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
