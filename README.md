# sequence

Sequence is intended to be a CLI for local development or a CI/CD server that can use the same pluggable container runtime(s) to produce artifacts--bridging the gap between local development and automation. Should support GitHub Actions as well as Concourse Resources. Written in Go.

## usage

see below for examples of a step, job or workflow's yaml

### step

```sh
# runs a step
$ sqnc run step step.yml
# runs a step of a job
$ sqnc run step -s=step job.yml
# runs a step of a job of a workflow
$ sqnc run step -j=job -s=step workflow.yml
```

### job

```sh
# runs a job
$ sqnc run job job.yml
# runs a job of a workflow
$ sqnc run job -j=job workflow.yml
```

### workflow

```sh
# runs a workflow
$ sqnc run workflow workflow.yml
```

## examples

### step

```yaml
# step.yml
image: golang:alpine
entrypoint:
  - go
cmd:
  - build
  - ./cmd/sqnc
```

### job

```yaml
# job.yml
steps:
  - image: golang:alpine
    entrypoint:
      - go
    cmd:
      - build
      - ./cmd/sqnc
  - image: golang:alpine
    entrypoint:
      - ./bin/sqnc
    cmd:
      - -v
```

### workflow

```yaml
# workflow.yml
jobs:
  example:
    steps:
      - image: golang:alpine
        entrypoint:
          - go
        cmd:
          - build
          - ./cmd/sqnc
      - image: golang:alpine
        entrypoint:
          - ./bin/sqnc
        cmd:
          - -v
```

## developing

TODO
