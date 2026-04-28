BINARY ?= dnshe-go
VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build run test clean

build:
	mkdir -p bin
	go build -trimpath -ldflags="$(LDFLAGS)" -o bin/$(BINARY) .

run:
	go run . -l 127.0.0.1:9876 -c data/config.json

test:
	go test ./...

clean:
	rm -rf bin
