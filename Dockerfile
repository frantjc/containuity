ARG base=alpine:3.15
ARG build=golang:1.18

FROM ${base} AS base

FROM ${build} AS build
ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/github.com/frantjc/sequence
COPY go.mod go.sum ./
RUN go mod download

FROM build AS bin
COPY . .
ARG version=0.0.0
ARG prerelease=
ARG commit=
ARG repository=frantjc/sequence
ARG tag=latest
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/bin ./cmd/sqncshim
RUN cp /usr/local/bin/sqncshim .
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/bin ./cmd/sqnc
RUN set -e; for pkg in $(go list ./...); do \
		go test -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/test/bin/$(basename $pkg).test -c $pkg; \
	done

FROM base AS sequence
COPY --from=bin /usr/local/bin /usr/local/bin
ENTRYPOINT ["sqnc"]

FROM sequence AS test
COPY --from=bin /usr/local/test/bin /usr/local/test/bin
RUN set -e; for test in /usr/local/test/bin/*.test; do \
		$test; \
	done

FROM sequence
