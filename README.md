# sequence

[![push](https://github.com/frantjc/sequence/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/sequence/actions)

<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/sequence/main/docs/demo.gif">
</p>

## developing

- `golang` is _required_ - version 1.18.x or above is required for generics
- `docker` is _required_ - version 20.10.x is tested
- `make` is _required_ - version 3.81 is tested
- `protoc` is _required_ if modifying the rpc API - version 3.19.x is tested
  - `protoc-gen-go`
- `buf` is _required_ if modifying the rpc API
  - `protoc-gen-connect-go`

```sh
# fmt
make vet
# install binary
make sqnc
# run rpc server
sqnc
# run workflows (usually requires github.token in ~/.sqnc/config)
sqnc run testdata/workflows/checkout_test_build_workflow.yml
```
