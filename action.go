package sequence

type Action struct {
	Name        string
	Author      string
	Description string
	Inputs      map[string]struct {
		Description        string
		Required           bool
		Default            string
		DeprecationMessage string
	}
	Outputs map[string]struct {
		Description string
	}
	Runs struct {
		Plugin     string
		Using      string
		Main       string
		Image      string
		Entrypoint string
		Args       []string
		Env        map[string]string
	}
}
