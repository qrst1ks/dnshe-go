BINARY ?= dnshe-go
VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build run test clean

build:
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) .

run:
	go run . -l 127.0.0.1:9999 -c data/config.json

test:
	go test ./...

clean:
	rm -rf bin $(BINARY)
