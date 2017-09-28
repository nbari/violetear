.PHONY: all bench build clean cover deps test

GO ?= go

all: build test

bench:
	${GO} test -run=^$$ -bench=.

deps:
	${GO} get github.com/nbari/violetear/middleware

build: deps
	${GO} build

clean:
	@rm -rf *.out

test: deps
	${GO} test -race

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
