package command

type Cmd interface {
	Execute() error
}
