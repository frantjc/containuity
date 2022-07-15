ARG base_image=alpine:3.15
ARG build_image=golang:1.18-alpine3.15

FROM ${base_image} AS base_image

FROM ${build_image} AS build_image
ENV CGO_ENABLED=0

FROM build_image AS build
WORKDIR $GOPATH/src/github.com/frantjc/sequence
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG version=0.0.0
ARG prerelease=
RUN go build -ldflags "-s -w" -o /usr/local/bin ./internal/cmd/sqnc-shim
RUN cp /usr/local/bin/sqnc-shim internal/shim/
RUN go build -ldflags "-s -w" -o /usr/local/bin ./internal/cmd/sqnc-runtime-docker
RUN go build -ldflags "-s -w" -o /usr/local/bin ./cmd/sqnc

FROM base_image AS sequence
COPY --from=build /usr/local/bin /usr/local/bin
RUN ln -s /usr/local/bin/sqnc-runtime-docker /etc/sqnc/plugins/sqnc-runtime-default
ENTRYPOINT ["sqnc"]

FROM sequence
