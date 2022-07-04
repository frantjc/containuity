package sequence

import "fmt"

const (
	ImageNode12        = Image("docker.io/library/node:12")
	ImageNode16        = Image("docker.io/library/node:16")
	DefaultRunnerImage = ImageNode12
)

type Image string

func (i Image) String() string {
	return string(i)
}

func (i Image) GoString() string {
	return fmt.Sprintf("sequence.Image(%s)", i)
}

func (i Image) GetRef() string {
	return i.String()
}
