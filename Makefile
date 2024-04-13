.PHONY: build run clean tidy

GOSELFSHIELD_SOURCES := $(shell find ./*.go ./cmd ./internal -type f -iname '*.go' ! -iname '*_test.go')

all: build

tidy:
	go mod tidy

build: goselfshield

goselfshield: $(GOSELFSHIELD_SOURCES) Makefile
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '-extldflags "-static"' -o goselfshield ./cmd/goselfshield

test: Makefile
	go test -v ./...

clean:
	rm -f ./goselfshield
