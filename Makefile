DOCKER ?= docker
GO ?= go

MODULE ?= github.com/frantjc/sequence
REPOSITORY ?= frantjc/sequence
TAG ?= latest

.PHONY: img
img:
	$(DOCKER) build -t $(REPOSITORY):$(TAG) --build-arg repository=$(REPOSITORY) --build-arg tag=$(TAG) .

.PHONY: bin
bin:
	$(GO) build -ldflags "-s -w -X $(MODULE).Repository=$(REPOSITORY) -X $(MODULE).Tag=$(TAG)" -o ./bin ./cmd/sqnc

.PHONY: all
all: img bin

.PHONY: test
test:
	$(DOCKER) build -t $(REPOSITORY):test --build-arg repository=$(REPOSITORY) --build-arg tag=test --target=test .

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...
