package sequence

const (
	ImageNode12        = image("docker.io/library/node:12")
	ImageNode16        = image("docker.io/library/node:16")
	DefaultRunnerImage = ImageNode12
)

type image string

func (i image) String() string {
	return string(i)
}

func (i image) GetRef() string {
	return i.String()
}
