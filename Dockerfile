ARG base=alpine:3.15
ARG build=golang:1.17-alpine3.15

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
ARG build=
ARG repository=frantjc/sequence
ARG tag=latest
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence/meta.Version=${version} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease}} -X github.com/frantjc/sequence/meta.Build=${build} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /usr/local/bin ./cmd/sqnc ./cmd/sqncd ./cmd/sqnctl
RUN set -e; for pkg in $(go list ./...); do \
		go test -ldflags "-s -w -X github.com/frantjc/sequence/meta.Version=${version} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease}} -X github.com/frantjc/sequence/meta.Build=${build} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /usr/local/test/bin/$(basename $pkg).test -c $pkg; \
	done

FROM base AS sequence
COPY --from=bin /usr/local/bin /usr/local/bin
ENTRYPOINT ["sqnc"]

FROM sequence AS test
COPY --from=sqnc /usr/local/test/bin /usr/local/test/bin
ARG SQNC_E2E
RUN set -e; for test in /usr/local/test/bin/*.test; do \
		$test; \
	done

FROM sequence
