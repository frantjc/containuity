BIN ?= /usr/local/bin

GO ?= go
GIT ?= git
DOCKER ?= docker
PROTOC ?= protoc

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

.DEFAULT: bin

bin: sqnc

bins binaries: sqnc sqncd

sqnc sqncd: shims
	@$(GO) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/$@
	@$(INSTALL) $(CURDIR)/bin/$@ $(BIN)

shims:
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/sqncshim
	@mv $(CURDIR)/bin/sqncshim $(CURDIR)/workflow/sqncshim
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/sqncshim-uses
	@mv $(CURDIR)/bin/sqncshim-uses $(CURDIR)/workflow/sqncshim-uses

placeholders:
	@cp $(CURDIR)/workflow/sqncshim.sh $(CURDIR)/workflow/sqncshim
	@cp $(CURDIR)/workflow/sqncshim.sh $(CURDIR)/workflow/sqncshim-uses

image img: 
	@$(DOCKER) build -t $(IMAGE) $(BUILD_ARGS) .

test: image
	@$(DOCKER) build -t $(REGISTRY)/$(REPOSITORY):test $(BUILD_ARGS) --target=test .

fmt lint pretty:
	@$(GO) fmt ./...

vet: fmt
	@$(GO) vet ./...

tidy: vet
	@$(GO) mod tidy

vendor: tidy
	@$(GO) mod vendor
	@$(GO) mod verify

clean: tidy placeholders
	@rm -rf bin/* vendor
	@$(DOCKER) system prune --volumes -a --filter label=sequence=true

protos:
	@$(PROTOC) $(PROTOC_ARGS) $(PROTOS)

coverage:
	@$(GO) test -v -cover -covermode=atomic -coverprofile=coverage.txt ./...

.PHONY: bin bins binaries sqnc sqncd shims placeholders image img test fmt lint pretty vet tidy vendor clean protos coverage
