# sequence

[![push](https://github.com/frantjc/sequence/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/sequence/actions)

<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/sequence/main/docs/demo.gif">
</p>

## developing

- `golang` is _required_ - version 1.16 or above is required for go mod to work
- `docker` is _required_ - version 20.10.x is tested
- `go mod` is _required_ for dependency management of golang packages
- `make` is _required_ - version 3.81 is tested
- `protoc` is _required_ if modifying the gRPC API - version 3.19.x is tested
    - `protoc-gen-go` - version 1.26
    - `protoc-gen-go-grpc` - version 1.1

```sh
# fmt
make vet
# install binary
make sqnc
# run gRPC server
sqnc
# run workflows (usually requires github.token in ~/.sqnc/config)
sqnc run testdata/workflows/checkout_test_build_workflow.yml
```
