GO ?= go

.PHONY: build
build:
	$(GO) build -o ./bin ./cmd/sqnc

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...
