package actions

type varsOpts struct {
	eopts []EnvOpt
	copts []CtxOpt
}

type VarsOpt func(v *varsOpts) error

func defaultVarsOpts() *varsOpts {
	return &varsOpts{}
}

func WithEnvOpts(opts ...EnvOpt) VarsOpt {
	return func(v *varsOpts) error {
		v.eopts = opts
		return nil
	}
}

func WithCtxOpts(opts ...CtxOpt) VarsOpt {
	return func(v *varsOpts) error {
		v.copts = opts
		return nil
	}
}
