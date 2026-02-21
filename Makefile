BINARY_NAME=dreadnotes

VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date +%FT%T%z)
PREFIX ?= /usr/local

LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

.PHONY: all build clean test run deps lint

all: build

build:
	@echo "Building ${BINARY_NAME} version ${VERSION}…"
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ../${BINARY_NAME}
	@echo "Build complete: bin/${BINARY_NAME}"

run: build
	@echo "Running…"
	./bin/${BINARY_NAME}

clean:
	@echo "Cleaning…"
	go clean
	rm -rf bin/

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

deps:
	go mod tidy
	go mod download

lint:
	golangci-lint run

install: build
	@echo "Installing to $(DESTDIR)$(PREFIX)/bin"
	install -D -m 755 bin/$(BINARY_NAME) $(DESTDIR)$(PREFIX)/bin/$(BINARY_NAME)

uninstall:
	@echo "Removing from $(DESTDIR)$(PREFIX)/bin"
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BINARY_NAME)
