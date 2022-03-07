ARG base=ubuntu:focal
ARG build=golang:1.17

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
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence/meta.Version=${version} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease} -X github.com/frantjc/sequence/meta.Build=${commit} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /usr/local/bin ./cmd/sequence
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence/meta.Version=${version} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease} -X github.com/frantjc/sequence/meta.Build=${commit} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /usr/local/bin ./cmd/sqnc
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence/meta.Version=${version} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease} -X github.com/frantjc/sequence/meta.Build=${commit} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /usr/local/bin ./cmd/sqncshim
RUN set -e; for pkg in $(go list ./...); do \
		go test -ldflags "-s -w -X github.com/frantjc/sequence/meta.Version=${version} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease} -X github.com/frantjc/sequence/meta.Build=${commit} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /usr/local/test/bin/$(basename $pkg).test -c $pkg; \
	done

FROM base AS sequence
COPY --from=bin /usr/local/bin /usr/local/bin
ENTRYPOINT ["sequence"]

FROM sequence AS test
COPY --from=bin /usr/local/test/bin /usr/local/test/bin
RUN set -e; for test in /usr/local/test/bin/*.test; do \
		$test; \
	done

FROM sequence
