package sequence

type Steppable interface {
	Steps() ([]Step, error)
}
