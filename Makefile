.PHONY: all build install clean

all: build

build:
	go build -v -o shorten main.go

run: build
	shorten

install:
	go get -v

clean:
	rm shorten
