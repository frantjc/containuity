DOCKER ?= docker
GO ?= go

MODULE ?= github.com/frantjc/sequence
REPOSITORY ?= frantjc/sequence
TAG ?= latest

.PHONY: test
test:
	$(DOCKER) build -t $(REPOSITORY):test --build-arg SQNC_E2E=yes --build-arg repository=$(REPOSITORY) --build-arg tag=test --target=test .

.PHONY: image
image: test
	$(DOCKER) build -t $(REPOSITORY):$(TAG) --build-arg repository=$(REPOSITORY) --build-arg tag=$(TAG) .

.PHONY: binaries
binaries: image
	$(DOCKER) run --rm --entrypoint sh -v `pwd`/bin:/assets $(REPOSITORY):$(TAG) -c "cp /usr/local/bin/* /assets"

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
	rm -rf bin/
	docker system prune --volumes -a --filter label=sequence=true
