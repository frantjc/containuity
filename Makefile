BIN ?= /usr/local/bin

GO ?= go
GO_LINUX_AMD64 ?= GOOS=linux GOARCH=amd64 $(GO)
GIT ?= git
DOCKER ?= docker
PROTOC ?= protoc
GOLANGCI-LINT ?= golangci-lint

VERSION ?= 0.0.0
PRERELEASE ?= alpha0

BRANCH ?= $(shell $(GIT) rev-parse --abbrev-ref HEAD 2>/dev/null)
COMMIT ?= $(shell $(GIT) rev-parse HEAD 2>/dev/null)
SHORT_SHA ?= $(shell $(GIT) rev-parse --short $(COMMIT))

REGISTRY ?= ghcr.io
REPOSITORY ?= frantjc/sequence
MODULE ?= github.com/$(REPOSITORY)
TAG ?= latest
IMAGE ?= $(REGISTRY)/$(REPOSITORY):$(TAG)

BUILD_ARGS ?= --build-arg version=$(VERSION) --build-arg prerelease=$(PRERELEASE)

PROTOS ?= $(shell find . -type f -name *.proto)
PROTOC_ARGS ?= --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative

INSTALL ?= sudo install

.DEFAULT: install

install: binaries
	@$(INSTALL) $(CURDIR)/bin/sqncd $(CURDIR)/bin/sqnc $(BIN)

bins binaries: sqnc sqncd

sqnc sqncd: shims
	@$(GO) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/$@

shims:
	@$(GO_LINUX_AMD64) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/sqncshim
	@cp $(CURDIR)/bin/sqncshim $(CURDIR)/workflow/sqncshim
	@$(GO_LINUX_AMD64) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/sqncshim-uses
	@cp $(CURDIR)/bin/sqncshim-uses $(CURDIR)/workflow/sqncshim-uses

placeholders:
	@cp $(CURDIR)/workflow/sqncshim.sh $(CURDIR)/workflow/sqncshim
	@cp $(CURDIR)/workflow/sqncshim.sh $(CURDIR)/workflow/sqncshim-uses

image img: 
	@$(DOCKER) build -t $(IMAGE) $(BUILD_ARGS) .

fmt vet test:
	@$(GO) $@ ./...

tidy vendor verify:
	@$(GO) mod $@

clean: tidy placeholders
	@rm -rf bin/* vendor
	@$(DOCKER) system prune --volumes -a --filter label=sequence=true

protos:
	@$(PROTOC) $(PROTOC_ARGS) $(PROTOS)

coverage:
	@$(GO) test -v -cover -covermode=atomic -coverprofile=coverage.txt ./...

lint:
	@$(GOLANGCI-LINT) run

tools:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	@echo 'Update your PATH so that the protoc compiler can find the plugins:'
	@echo '$$ export PATH=PATH:$(shell $(GO) env GOPATH)/bin"'

.PHONY: install bins binaries sqnc sqncd shims placeholders image img fmt vet test tidy vendor verify clean protos coverage lint tools
