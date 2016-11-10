#!/usr/bin/make -f

GLIDE=$(shell which glide)
PACKAGES=./engine ./ui
SRC=$(shell find . -maxdepth 2 -type f -name '*.go' -a ! -name '*_test.go')
SRC_TEST=$(shell find . -maxdepth 2 -type f -name '*_test.go')

all: dd-logstats deps

deps:
	glide install

test: $(SRC) $(SRC_TEST)
	go test -v $(PACKAGES)

dd-logstats: deps $(SRC)
	go build -o $@

	
.PHONY=all test deps
