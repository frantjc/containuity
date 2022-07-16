# sequence

[![push](https://github.com/frantjc/sequence/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/sequence/actions)

Run _sequential_ containerized workloads on the same volume using tools from each container along the way.

<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/sequence/main/docs/demo.gif">
</p>

## summary

Sequence is, first and foremost, a library for running _sequential_ containerized workloads on the same volume to produce some result. To achieve this, it builds upon some existing technologies:

- [x] Borrow [GitHub Action's Workflow syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions) and expand upon it (e.g. by allowing each step to designate what image it should run inside). This should allow Sequence to utilize useful GitHub Actions, which are GitHub repositories which can be concisely referenced to execute complicated tasks (e.g. [downloading and installing Go](https://github.com/actions/setup-go))
- [x] Use a pluggable container [`Runtime`](runtime/runtime.go) whose default implementation is [Docker](https://docs.docker.com/get-started/) to run each containerized task.
- [ ] Take advantage of [Concourse Resources](https://concourse-ci.org/resources.html) to additionally expand the functionality of what a single step can do

Sequence aims to have tools built from this library to unify the development and continuous integration (CI) experiences:

- [ ] `sqnctl` CLI to run workflows against local changes before pushing them to be executed by CI
- [ ] `sqncd` RPC daemon that can be connected to remotely to run workloads

## developing

- `make` is _required_ - version 3.81 is tested
- `golang` is _required_ - version 1.18.x or above is required for [generics](https://go.dev/doc/tutorial/generics)
- `docker` is _required_ - version 20.10.x is tested
- [`buf`](https://github.com/bufbuild/buf) is _required if_ modifying proto - version 1.4.x is tested
- [`upx`](https://github.com/upx/upx) is _required_ for compressing [`sqnc-shim`](internal/shim/sqnc-shim) on generate
- [`protoc`](https://grpc.io/docs/protoc-installation/) is _required if_ modifying proto - version 3.19.x is tested
  - [`protoc-gen-go`](https://developers.google.com/protocol-buffers/docs/reference/go-generated) - version 1.26.x is tested
  - (hopefully) temporarily, [`protoc-gen-sqnc-go`](internal/cmd/protoc-gen-sqnc/main.go)

The latter two of these can be installed by:

```sh
make tools
```

### test

Create a `.env` that looks like [`.env.example`](.env.example) but with a _real_ GitHub token, and:

```sh
make test # go test ./...
```

### lint

Format `.go` code.

```sh
make fmt    # go fmt ./...
make lint   # golangci-lint run
```

### generate

Generate `.go` code from `.proto` code.

```sh
make generate # buf generate .
```
