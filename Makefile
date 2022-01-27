DOCKER ?= docker
GO ?= go

MODULE ?= github.com/frantjc/sequence
REPOSITORY ?= frantjc/sequence
TAG ?= latest

SQNC_E2E ?=

.PHONY: test
test:
	$(DOCKER) build -t $(REPOSITORY):test --build-arg SQNC_E2E=$(SQNC_E2E) --build-arg repository=$(REPOSITORY) --build-arg tag=test --target=test .

.PHONY: image
image: test
	$(DOCKER) build -t $(REPOSITORY):$(TAG) --build-arg repository=$(REPOSITORY) --build-arg tag=$(TAG) .

.PHONY: binaries
binaries: image
	$(DOCKER) build -t $(REPOSITORY):bin --build-arg repository=$(REPOSITORY) --build-arg tag=$(TAG) --target=bin .
	$(DOCKER) run --rm -v `pwd`/bin:/assets $(REPOSITORY):bin sh -c "cp /usr/local/bin/* /assets"
	chmod +x bin/*

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...

.PHONY: all
all: vet binaries

.PHONY: clean
clean:
	rm -rf bin/*
	docker system prune --volumes -a --filter label=sequence=true
