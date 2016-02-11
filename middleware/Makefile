.PHONY: all deps test build cover

GO ?= go

all: build test

build:
	${GO} build

test: deps
test:
	${GO} test -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
