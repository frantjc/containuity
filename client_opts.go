package sequence

type clientOpts struct{}

type ClientOpt func(co *clientOpts) error
