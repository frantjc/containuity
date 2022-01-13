package actions

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Reference interface {
	Owner() string
	Repository() string
	Path() string
	Version() string
	String() string
}

type reference struct {
	owner      string
	repository string
	path       string
	version    string
}

func (r *reference) Owner() string {
	return r.owner
}

func (r *reference) Repository() string {
	return r.repository
}

func (r *reference) Path() string {
	return r.path
}

func (r *reference) Version() string {
	return r.version
}

func (r *reference) String() string {
	s := fullRepository(r)
	if r.path != "" {
		s = fmt.Sprintf("%s/%s", s, r.path)
	}
	if r.version != "" {
		s = fmt.Sprintf("%s@%s", s, r.version)
	}
	return s
}

func ParseReference(actionRef string) (Reference, error) {
	r := &reference{}

	spl1 := strings.Split(actionRef, "@")
	switch len(spl1) {
	case 2:
		r.version = spl1[1]
	case 1:
	default:
		return r, fmt.Errorf("unable to parse action ref %s", actionRef)
	}

	spl2 := strings.Split(spl1[0], "/")
	switch len(spl2) {
	case 0, 1:
		return r, fmt.Errorf("unable to parse action ref %s", actionRef)
	default:
		r.owner = spl2[0]
		r.repository = spl2[1]
		if len(spl2) > 2 {
			r.path = filepath.Join(spl2[2:]...)
		}
	}

	return r, nil
}

func fullRepository(r Reference) string {
	return fmt.Sprintf("%s/%s", r.Owner(), r.Repository())
}
