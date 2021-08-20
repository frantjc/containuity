# containuity

In large part for education and enjoyment purposes, Containuity will be an attempt at an open-source, container-native CI tool that allows for a continuous user experience between developing locally and normal CI processes.

That is, it should allow developers to overwrite inputs to a task (e.g. source code in a Git repository) with an input on their machine without pushing that input to a remote destination, while still allowing the use of a remote destination as an input to a task.

Additionally, it should allow developers to tell Containuity to send task outputs to their local machine rather than to some remote destination, while still allowing the use of a remote destination as the output destination of a task.

## architecture

- everything external to a pipeline that gets pulled in should be modelable as a "resource" a la [Concourse](https://github.com/concourse/concourse)
- external database not required; Containuity can fulfil its own database requirements by running one in a container at runtime if no external database is supplied
- if Containuity is fulfilling its own database requirements, agents should maintain replicas of the controller's database on them
- Containuity's controller can execute tasks, but should delegate tasks to agents if possible
- container runtime should be pluggable but will rely on [containerd](https://github.com/containerd/containerd) primarily
