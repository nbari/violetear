.PHONY: all test clean build install

GO ?= go
GOPATH := $(CURDIR)/_vendor:$(GOPATH)

all: install test

build:
	@go build

clean:
	@go clean

install:
	@go get

test: install
	@go test -v
