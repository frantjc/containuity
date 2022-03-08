# sequence

## developing

- `golang` is _required_ - version 1.16 or above is required for go mod to work
- `docker` is _required_ - version 20.10.x is tested
- `go mod` is _required_ for dependency management of golang packages
- `make` is _required_ - version 3.81 is tested
- `protoc` is _required_ - version 3.19.x is tested
    - `protoc-gen-go` - version 1.26
    - `protoc-gen-go-grpc` - version 1.1

```sh
# install binaries
make
# build frantjc/sequence image
make image
# fmt
make pretty
```
