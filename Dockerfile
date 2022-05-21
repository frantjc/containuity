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
ARG commit=
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/bin ./cmd/sqncshim
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/bin ./cmd/sqncshim-uses
RUN cp /usr/local/bin/sqncshim ./workflow
RUN cp /usr/local/bin/sqncshim-uses ./workflow
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/bin ./cmd/sqnc
RUN go build -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/bin ./cmd/sqncd
RUN set -e; for pkg in $(go list ./...); do \
		go test -ldflags "-s -w -X github.com/frantjc/sequence.Version=${version} -X github.com/frantjc/sequence.Prerelease=${prerelease} -X github.com/frantjc/sequence.Build=${commit}" -o /usr/local/test/bin/$(basename $pkg).test -c $pkg; \
	done

FROM base_image AS sequence
COPY --from=build /usr/local/bin /usr/local/bin
ENTRYPOINT ["sqncd"]

FROM sequence AS test
COPY --from=build /usr/local/test/bin /usr/local/test/bin
RUN set -e; for test in /usr/local/test/bin/*.test; do \
		$test; \
	done

FROM sequence
