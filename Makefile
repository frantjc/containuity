BIN ?= /usr/local/bin
DOTENV ?= .env

GO ?= go
GIT ?= git
DOCKER ?= docker
GOLANGCI-LINT ?= golangci-lint
BUF ?= buf

VERSION ?= 0.0.0
PRERELEASE ?= dev0

BRANCH ?= $(shell $(GIT) rev-parse --abbrev-ref HEAD 2>/dev/null)
COMMIT ?= $(shell $(GIT) rev-parse HEAD 2>/dev/null)
SHORT_SHA ?= $(shell $(GIT) rev-parse --short $(COMMIT))

REGISTRY ?= ghcr.io
REPOSITORY ?= frantjc/sequence
MODULE ?= github.com/$(REPOSITORY)
TAG ?= latest
IMAGE ?= $(REGISTRY)/$(REPOSITORY):$(TAG)

BUILD_ARGS ?= --build-arg version=$(VERSION) --build-arg prerelease=$(PRERELEASE)

GITHUB_TOKEN ?= $(shell source $(DOTENV) && echo $$GITHUB_TOKEN)

INSTALL ?= sudo install

.DEFAULT: install

install: binaries
	@$(INSTALL) $(CURDIR)/bin/sqncd $(CURDIR)/bin/sqnc $(BIN)

bins binaries: sqnc sqncd sqnctl

sqnc sqncd sqnctl:
# @$(GO) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/$@

image img: 
	@$(DOCKER) build -t $(IMAGE) $(BUILD_ARGS) .

test:
	@GITHUB_TOKEN=$(GITHUB_TOKEN) $(GO) $@ -v ./...

fmt vet generate:
	@$(GO) $@ ./...

tidy vendor verify download:
	@$(GO) mod $@

clean: tidy
	@rm -rf bin/* vendor
	@$(DOCKER) system prune --volumes -a --filter label=sequence=true

lint:
	@$(GOLANGCI-LINT) run

tools:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@$(GO) install ./internal/cmd/protoc-gen-sqnc
	@echo 'Update your PATH so that the protoc compiler can find the plugins:'
	@echo '$$ export PATH=$$PATH:$(shell $(GO) env GOPATH)/bin"'

.PHONY: \
	install bins binaries sqnc sqncd \
	shim shims sqnc-shim \
	image img \
	generate test fmt vet \
	tidy vendor verify \
	clean \
	lint tools
