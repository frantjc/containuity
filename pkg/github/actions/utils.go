package actions

func coalesce(a ...string) string {
	for _, s := range a {
		if s != "" {
			return s
		}
	}
	return ""
}
