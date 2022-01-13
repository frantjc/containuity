package runtime

type Pull struct{}

type PullOpt func(po *Pull) error

func NewPull(opts ...PullOpt) (*Pull, error) {
	p := &Pull{}
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
