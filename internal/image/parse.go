package image

import (
	"github.com/google/go-containerregistry/pkg/name"
)

func ParseRef(ref string) (string, error) {
	pref, err := name.ParseReference(ref)
	if err != nil {
		return "", err
	}

	return pref.Name(), nil
}
