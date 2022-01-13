# architecture

## daemon

what sequence does:

* takes a Step, Job or Workflow and runs it in a series of sequential containers
* takes a Step and attaches to it
* takes an GitHub Action ref and a path, clones the Action repository into that path, and the emits the action as a Step to stdout

what sequence's daemon needs to do:

* take a Step, Job or Workflow and rus it in a series of sequential containers
* take a Step and attach to it

what sequence's daemon does not need to do:

* take an GitHub Action ref and a path, clones the Action repository into that path, and the emits the action as a Step to stdout
