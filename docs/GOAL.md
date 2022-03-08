# goal

sequence aims to be a tool that can be used to run _sequential_ containerized workloads that operate on the same volume(s) using whatever tools are native to each container along the way à la [concourse](https://concourse-ci.org)

by virtue of containerization, sequence should be able to utilize a pluggable container runtime. By running workloads using the [docker](https://docker.com) it should be easy to use as a local development tool. By running workloads using [kubernetes](https://kubernetes.io/), sequence should be able to scale with some ease

by virute of [golang](https://go.dev/), sequence should be a viable tool on any operating system

because of these things, sequence aims to be both a cli tool as well as something of a continuous integration server, bridging the gap that often exists between local development and continuous integration by allowing pipelines to be runnable manually on a developer's machine as well as automatically on some server

by virtue of tools with similar goals that have come before it, sequence should be able to take advantage of concourse [resources](https://concourse-ci.org/resources.html) as well as github [actions](https://docs.github.com/en/actions/learn-github-actions/understanding-github-actions#actions) to avoid having to implement its own versions of common tasks, such as checking out a git repository

## first draft

as a first draft, sequence should strictly be a cli for running github actions workflows locally, as that has the most potential value, as github actions is presumably the most widely used ci system today. support for concourse resources as well as server functionality can be added later

## architecture

### workflow

[example](testdata/workflow.yml)

a sequence workflow aims to be a superset of a github actions [workflow](https://docs.github.com/en/actions/learn-github-actions/understanding-github-actions#create-an-example-workflow), consisting of "jobs" of containerized workloads as well as services that are exposed to those workloads

### job

[example](testdata/job.yml)

a sequence job aims to be a superset of a github actions job

### step

[example](testdata/step.yml)

a sequence step aims to be a superset of a github actions step as well as a concourse [step](https://concourse-ci.org/steps.html)

in an ideal world, a step should consist of a single containerized workload. however, a github action isn't actually ran in a single environment; rather, a step specifying [uses](https://docs.github.com/en/actions/learn-github-actions/workflow-syntax-for-github-actions#jobsjob_idstepsuses) it gathers some information about the action on the host machine, perhaps builds the container for the action in another and exeutes the action. where github has the benefit of ease of control over the environment in which they are running such steps, an easy-to-use dev tool such as this does not...

### special step

sequence aims to allow such steps to exist via "special" steps that emit information about what the subsequent step should be in the form of json to stdout. luckily, sequence would have had to do something like this eventually, anyways, since concourse [resource types](https://concourse-ci.org/implementing-resource-types.html) already transmit information this way

if sequence is running on linux, it should be able to mount _itself_ into an arbitrary container to execute such steps. however, if sequence is not running on linux, a linux container wouldn't be able to execute, say, a darwin build of sequence, and so will have to pull the sequence image itself to accomplish the same

#### uses

sequence's "uses" plugin should clone a given action to a given path, parse that action's action.yml file, and then emit the "step" version of that action to stdout, letting sequence know what the following step should be

#### github

some older, "core" github actions maintained by github (e.g. [actions/checkout@v1](https://github.com/actions/checkout/blob/v1/action.yml#L23)) utilize "plugins" that live on github actions runners themselves, and so will have to be implemented by sequence itself if support for such actions is desired

## daemon

sequence should be able to run as a daemon--a gRPC client/server combination--for speed purposes. additionally, sequence should be able to run daemonless--just making calls the the daemon(s) of the container runtime of choice.

## cli

### `sqncd`

run the sequence gRPC server

### `sqnc run`

run a step

```sh
sqnc run step [FLAGS...] STEP | -s=STEP_ID JOB | -j=JOB_ID -s=STEP_ID WORKFLOW
```

run the steps of a job

```sh
sqnc run job [FLAGS...] JOB | -j=JOB_ID WORKFLOW
```

run the jobs of a workflow

```sh
sqnc run workflow [FLAGS...] WORKFLOW
```

### `sqnc debug`

get a shell inside of the container of a would-be step

```sh
sqnc debug step [FLAGS...] STEP | -s=STEP_ID JOB | -j=JOB_ID -s=STEP_ID WORKFLOW
```

### `sqncshim plugin`

clones a given action and emits the step version of that action to stdout

```sh
sqncshim plugin uses [FLAGS...] ACTION [PATH]
```

runs sequence's own implementation of a github actions plugin

```sh
sqncshim plugin [FLAGS...] PLUGIN
```

## build

a single binary that's able to run as a gRPC daemon, daemonless-ly, and with interchangable container runtime backends sounds pretty big. it might make sense to build sequence into several different binaries with different combinations of functionality. first thought:

* `sqncshim` -- sequence shim, contains plugins

* `sqncd` -- gRPC server itself

* `sqnctl` -- cli to interact with gRPC server
