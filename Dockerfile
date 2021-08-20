ARG alpine_image=alpine:edge
ARG default_image=debian:bookworm
ARG build_image=golang:latest

FROM ${alpine_image} AS alpine

FROM ${default_image} AS default
ENV DEBIAN_FRONTEND noninteractive

FROM ${build_image} AS build
ENV CGO_ENABLED 0
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY worker/ worker/
COPY *.go ./
ARG ldflags
RUN go build -ldflags "${ldflags}" -o /assets/worker ./worker/cmd/ \
    && go build -ldflags "${ldflags}" -o /assets/geocloud ./cmd/

FROM default AS containerd
ARG containerd=https://github.com/containerd/containerd/releases/download/v1.5.5/containerd-1.5.5-linux-amd64.tar.gz
# when its src is a remote .tgz, ADD does not unpack the tarball
# when its src is a local .tgz, ADD unpacks the tarball
ADD ${containerd} /tmp/
# this conditional handles that difference in ADD's functionality between remote and local
RUN TGZ=/tmp/$(basename ${containerd}); \
    if [ -f $TGZ ]; then \
        tar \
            -C /tmp/ \
            -xzf \
            $TGZ \
        && rm $TGZ; \
    fi; \
    mkdir /assets/ \
    && mv /tmp/bin/* /assets/ \
    && chmod +x /assets/*

FROM default AS runc
ARG runc=https://github.com/opencontainers/runc/releases/download/v1.0.1/runc.amd64
ADD ${runc} /assets/runc
RUN chmod +x /assets/runc

FROM alpine AS containuity_alpine
RUN apk update \
    && apk add --no-cache \
        ca-certificates \
        tini \
        pigz
COPY --from=build /assets/worker /usr/local/containuity/bin/containuity
COPY --from=containerd /assets/ /usr/local/containuity/bin/
COPY --from=runc /assets/ /usr/local/containuity/bin/
VOLUME /var/lib/containuity/containerd/
ENV PATH=/usr/local/containuity/bin:$PATH
ENTRYPOINT ["tini", "containuity"]
COPY --from=build /assets/containuity /usr/local/containuity/bin/containuity

FROM default AS containuity
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        dumb-init \
        pigz \
    && rm -rf /var/lib/apt/lists/*
COPY --from=build /assets/worker /usr/local/containuity/bin/containuity
COPY --from=containerd /assets/ /usr/local/containuity/bin/
COPY --from=runc /assets/ /usr/local/containuity/bin/
VOLUME /var/lib/containuity/containerd/
ENV PATH=/usr/local/containuity/bin:$PATH
ENTRYPOINT ["dumb-init", "containuity"]
COPY --from=build /assets/containuity /usr/local/containuity/bin/containuity

FROM containuity
