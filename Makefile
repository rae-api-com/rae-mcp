BINARY=rae-mcp
BINDIR=bin
PREFIX?=/usr/local/bin

.PHONY: all build install clean run fmt

all: build

build:
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINDIR)/$(BINARY) .

install: build
	install -c $(BINDIR)/$(BINARY) $(PREFIX)/$(BINARY)

clean:
	rm -rf $(BINDIR)

run:
	go run *.go

fmt:
	go fmt ./...
	golines -w .
	goimports -w .
