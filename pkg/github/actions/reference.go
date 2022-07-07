package actions

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Reference struct {
	Owner      string
	Repository string
	Path       string
	Version    string
}

func (r *Reference) String() string {
	s := fullRepository(r)
	if r.Path != "" {
		s = fmt.Sprintf("%s/%s", s, r.Path)
	}
	if r.Version != "" {
		s = fmt.Sprintf("%s@%s", s, r.Version)
	}
	return s
}

func ParseReference(actionRef string) (*Reference, error) {
	r := &Reference{}

	spl1 := strings.Split(actionRef, "@")
	switch len(spl1) {
	case 2:
		r.Version = spl1[1]
	case 1:
	default:
		return r, fmt.Errorf("unable to parse action ref %s", actionRef)
	}

	spl2 := strings.Split(spl1[0], "/")
	switch len(spl2) {
	case 0, 1:
		return r, fmt.Errorf("unable to parse action ref %s", actionRef)
	default:
		r.Owner = spl2[0]
		r.Repository = spl2[1]
		if len(spl2) > 2 {
			r.Path = filepath.Join(spl2[2:]...)
		}
	}

	return r, nil
}

func fullRepository(r *Reference) string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Repository)
}
