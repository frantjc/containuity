BIN ?= /usr/local/bin

GO ?= go
GOOS ?= $(shell $(GO) env GOOS)

GIT ?= git
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null)

REGISTRY ?= docker.io
REPOSITORY ?= frantjc/sequence
MODULE ?= github.com/$(REPOSITORY)
TAG ?= latest

E2E ?= yes

DOCKER ?= docker

BUILD_ARGS ?= --build-arg repository=$(REPOSITORY) --build-arg SQNC_E2E=$(E2E) --build-arg tag=$(TAG) --build-arg commit=$(COMMIT)

.DEFAULT: sqnc

sqnc:
	$(GO) build -ldflags "-s -w -X github.com/frantjc/sequence/meta.Build=$(COMMIT) -X github.com/frantjc/sequence/meta.Repository=$(REPOSITORY) -X github.com/frantjc/sequence/meta.Tag=$(TAG)" -o ./bin $(CURDIR)/cmd/$@
	sudo install $(CURDIR)/bin/$@ $(BIN)

image img: 
	$(DOCKER) build -t $(REGISTRY)/$(REPOSITORY):$(TAG) $(BUILD_ARGS) .

test: image
	$(DOCKER) build -t $(REGISTRY)/$(REPOSITORY):test $(BUILD_ARGS) --target=test .

bin bins binaries: sqnc

fmt lint pretty:
	$(GO) fmt ./...

vet: fmt
	$(GO) vet ./...

vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify

clean:
	rm -rf bin/* vendor
	docker system prune --volumes -a --filter label=sequence=true

.PHONY: sqnc image img test bin bins binaries fmt lint pretty vet vendor clean
