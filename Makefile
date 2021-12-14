DOCKER ?= docker
GO ?= go

REPOSITORY ?= frantjc/sequence
TAG ?= latest

.PHONY: build
build:
	$(DOCKER) build -t $(REPOSITORY):$(TAG) --build-arg repository=$(REPOSITORY) --build-arg tag=$(TAG) .
	$(GO) build -o ./bin ./cmd/sqnc

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...
