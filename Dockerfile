ARG base_image=debian
ARG build_image=golang:1.18

FROM ${base_image} AS base_image

FROM ${build_image} AS build_image

FROM build_image AS build
WORKDIR $GOPATH/src/github.com/frantjc/sequence
ARG version=0.0.0
ARG prerelease=
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o $GOPATH/bin ./cmd/sqnc
RUN go build -buildmode=plugin -ldflags "-s -w" -o $GOPATH/bin/sqnc-runtime-docker.so ./internal/cmd/sqnc-runtime-docker

FROM base_image AS sequence
RUN mkdir -p /etc/sqnc/plugins
COPY --from=build /go/bin /usr/local/bin
RUN mv /usr/local/bin/*.so /etc/sqnc/plugins
ENTRYPOINT ["sqnc"]

FROM sequence
