NAME = bcrypt-tool

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
CRYPTO_VERSION ?= $(shell grep 'golang.org/x/crypto' go.mod | awk '{print $$NF}')

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION) -X main.cryptoVersion=$(CRYPTO_VERSION)" -o output/$(NAME)

.PHONY: clean
clean:
	rm -rf dist output/$(NAME)

.PHONY: test
test:
	go test -race ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: release
release:
	envy exec gh-release goreleaser release --clean
	$(MAKE) clean

default: build
