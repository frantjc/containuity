package command

type fileDescriptor interface {
	Fd() uintptr
}
