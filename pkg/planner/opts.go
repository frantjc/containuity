package planner

type planOpts struct {
	path string
}

type PlanOpt func(*planOpts) error

func WithPath(path string) PlanOpt {
	return func(po *planOpts) error {
		po.path = path
		return nil
	}
}
