.PHONY: all build clean

GO=go
BIN=truenascsi

all: clean build

build:/
	$(GO) build -o bin/$(BIN) cmd/truenascsi/*.go

clean:
	rm -rf bin/$(BIN)