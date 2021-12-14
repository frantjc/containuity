ARG base=alpine:3.15
ARG build=golang:1.17-alpine3.15

FROM ${base} AS base

FROM ${build} AS build
ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/github.com/frantjc/sequence
COPY go.mod go.sum ./
RUN go mod download

FROM build AS sqnc
COPY . .
ARG version=0.0.0
ARG revision=
ARG repository=frantjc/sequence
ARG tag=latest
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${verision} -X github.com/frantjc/sequence.Revision=${revision} -X github.com/frantjc/sequence.Repository=${repository} -X github.com/frantjc/sequence.Tag=${tag} -X github.com/frantjc/sequence.Digest=${digest}" -o /assets/sqnc ./cmd/sqnc/

FROM base AS sequence
COPY --from=sqnc /assets/sqnc /usr/local/bin/sqnc
ENTRYPOINT ["sqnc"]
