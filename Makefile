.PHONY: all build install lint clean

all: build

build:
	go build -v -o shorten main.go

install:
	go get -v

lint:
	${GOPATH}/bin/golint . cmd markup

clean:
	rm shorten
