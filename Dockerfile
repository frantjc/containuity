ARG base_image=alpine:3.15
ARG build_image=golang:1.18-alpine3.15

FROM ${base_image} AS base_image

FROM ${build_image} AS build_image
ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/github.com/frantjc/sequence
COPY go.mod go.sum ./
RUN go mod download

FROM build_image AS build
COPY . .
ARG version=0.0.0
ARG prerelease=
RUN go build -ldflags "-s -w" -o /usr/local/bin ./internal/cmd/shim/source
RUN go build -ldflags "-s -w" -o /usr/local/bin ./internal/cmd/shim/uses
RUN cp /usr/local/bin/source ./workflow
RUN cp /usr/local/bin/uses ./workflow
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease}" -o /usr/local/bin ./cmd/sqnc
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease}" -o /usr/local/bin ./cmd/sqncd

FROM base_image AS sequence
COPY --from=build /usr/local/bin /usr/local/bin
ENTRYPOINT ["sqncd"]

FROM sequence
