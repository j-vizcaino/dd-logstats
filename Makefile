#!/usr/bin/make -f

GLIDE=$(shell which glide)
SRC=$(shell ls -1 *.go | grep -v '_test.go')

all: dd-logstats deps

deps:
	glide install

test:
	go test -v

dd-logstats: deps $(SRC)
	go build -o $@

	
.PHONY=all test deps
