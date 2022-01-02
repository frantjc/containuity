package actions

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Uses struct {
	Owner      string
	Repository string
	Path       string
	Version    string
}

func (u *Uses) Repo() string {
	return fmt.Sprintf("%s/%s", u.Owner, u.Repository)
}

func (u *Uses) String() string {
	s := u.Repo()
	if u.Path != "" {
		s = fmt.Sprintf("%s/%s", s, u.Path)
	}
	if u.Version != "" {
		s = fmt.Sprintf("%s@%s", s, u.Version)
	}
	return s
}

func Parse(uses string) (*Uses, error) {
	u := &Uses{}

	spl1 := strings.Split(uses, "@")
	switch len(spl1) {
	case 2:
		u.Version = spl1[1]
	case 1:
	default:
		return u, fmt.Errorf("unable to parse action %s", uses)
	}

	spl2 := strings.Split(spl1[0], "/")
	switch len(spl2) {
	case 0, 1:
		return u, fmt.Errorf("unable to parse action %s", uses)
	default:
		u.Owner = spl2[0]
		u.Repository = spl2[1]
		if len(spl2) > 2 {
			u.Path = filepath.Join(spl2[2:]...)
		}
	}

	return u, nil
}
