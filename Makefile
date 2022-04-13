BIN ?= /usr/local/bin

GO ?= go
GIT ?= git
DOCKER ?= docker
PROTOC ?= protoc

BRANCH ?= $(shell $(GIT) rev-parse --abbrev-ref HEAD 2>/dev/null)
COMMIT ?= $(shell $(GIT) rev-parse HEAD 2>/dev/null)
SHORT_SHA ?= $(shell $(GIT) rev-parse --short $(COMMIT))

REGISTRY ?= ghcr.io
REPOSITORY ?= frantjc/sequence
MODULE ?= github.com/$(REPOSITORY)
TAG ?= latest
IMAGE ?= $(REGISTRY)/$(REPOSITORY):$(TAG)

BUILD_ARGS ?= --build-arg repository=$(REPOSITORY) --build-arg tag=$(TAG) --build-arg commit=$(SHORT_SHA)

PROTOS ?= $(shell find . -type f -name *.proto)
PROTOC_ARGS ?= --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative

INSTALL ?= sudo install

.DEFAULT: bin

bin: sqnc

bins binaries: sqnc sqncd sqnctl

sqnc sqncd sqnctl: shims
	$(GO) build -ldflags "-s -w -X github.com/frantjc/sequence.Build=$(SHORT_SHA) -X github.com/frantjc/sequence.Tag=$(TAG)" -o $(CURDIR)/bin $(CURDIR)/cmd/$@
	$(INSTALL) $(CURDIR)/bin/$@ $(BIN)

shims:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o $(CURDIR)/bin $(CURDIR)/cmd/sqncshim
	mv $(CURDIR)/bin/sqncshim $(CURDIR)/workflow/sqncshim
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o $(CURDIR)/bin $(CURDIR)/cmd/sqncshim-uses
	mv $(CURDIR)/bin/sqncshim-uses $(CURDIR)/workflow/sqncshim-uses

image img: 
	$(DOCKER) build -t $(IMAGE) $(BUILD_ARGS) .

test: image
	$(DOCKER) build -t $(REGISTRY)/$(REPOSITORY):test $(BUILD_ARGS) --target=test .

fmt lint pretty:
	$(GO) fmt ./...

vet: fmt
	$(GO) vet ./...

tidy: vet
	$(GO) mod tidy

vendor: tidy
	$(GO) mod vendor
	$(GO) mod verify

clean: tidy
	rm -rf bin/* vendor
	$(DOCKER) system prune --volumes -a --filter label=sequence=true

protos:
	$(PROTOC) $(PROTOC_ARGS) $(PROTOS)

coverage:
	$(GO) test -v -cover -covermode=atomic -coverprofile=coverage.txt ./...

.PHONY: bin bins binaries sqnc sqncd sqnctl sqncshim shim image img test fmt lint pretty vet tidy vendor clean protos coverage
