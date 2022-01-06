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
ARG major=0
ARG minor=0
ARG patch=0
ARG prerelease=
ARG build=
ARG repository=frantjc/sequence
ARG tag=latest
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence/meta.Major=${major} -X github.com/frantjc/sequence/meta.Minor=${minor} -X github.com/frantjc/sequence/meta.Patch=${patch} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease}} -X github.com/frantjc/sequence/meta.Build=${build} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /assets/sqnc ./cmd/sqnc/
RUN set -e; for pkg in $(go list ./...); do \
		go test -ldflags "-s -w -X github.com/frantjc/sequence/meta.Major=${major} -X github.com/frantjc/sequence/meta.Minor=${minor} -X github.com/frantjc/sequence/meta.Patch=${patch} -X github.com/frantjc/sequence/meta.Prerelease=${prerelease}} -X github.com/frantjc/sequence/meta.Build=${build} -X github.com/frantjc/sequence/meta.Repository=${repository} -X github.com/frantjc/sequence/meta.Tag=${tag} -X github.com/frantjc/sequence/meta.Digest=${digest}" -o /assets/tests/$(basename $pkg).test -c $pkg; \
	done

FROM base AS sequence
COPY --from=sqnc /assets/sqnc /usr/local/bin
ENTRYPOINT ["sqnc"]

FROM sequence AS test
COPY --from=sqnc /assets/tests /usr/local/bin/tests
ARG SQNC_TEST_ACTION
RUN set -e; for test in /usr/local/bin/tests/*.test; do \
		$test; \
	done

FROM sequence
