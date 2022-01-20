package meta

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/opencontainers/go-digest"
)

var (
	Registry = "docker.io"

	Repository = "frantjc/sequence"
)

var (
	Tag = "latest"

	Digest = ""

	GoDigest digest.Digest

	ImageRef name.Reference
)

var (
	VirtualTag = "virtual"

	VirtualDigest = ""

	VirtualGoDigest digest.Digest

	VirtualImageRef name.Reference
)

func init() {
	if Repository == "" {
		panic(fmt.Sprintf("%s/meta.Repository must not be empty", Module))
	}
}

func init() {
	var (
		imageRef = Repository
		err      error
	)

	if Registry != "" {
		imageRef = fmt.Sprintf("%s/%s", Registry, imageRef)
	}

	if Tag != "" {
		imageRef = fmt.Sprintf("%s:%s", imageRef, Tag)
	}

	if Digest != "" {
		GoDigest = digest.FromString(Digest)
		err := GoDigest.Validate()
		if err != nil {
			panic(fmt.Sprintf("%s/meta.Digest is invalid: %s", Module, err.Error()))
		}
		imageRef = fmt.Sprintf("%s@%s", imageRef, GoDigest.String())
	}

	ImageRef, err = name.ParseReference(imageRef)
	if err != nil {
		panic(fmt.Sprintf("%s/meta.ImageRef is invalid: %s", Module, err.Error()))
	}
}

func Image() string {
	return ImageRef.Name()
}

func init() {
	var (
		virtualImageRef = Repository
		err             error
	)

	if Registry != "" {
		virtualImageRef = fmt.Sprintf("%s/%s", Registry, virtualImageRef)
	}

	if Tag != "" {
		virtualImageRef = fmt.Sprintf("%s:%s", virtualImageRef, VirtualTag)
	}

	if Digest != "" {
		VirtualGoDigest = digest.FromString(VirtualDigest)
		err := VirtualGoDigest.Validate()
		if err != nil {
			panic(fmt.Sprintf("%s/meta.VirtualDigest is invalid: %s", Module, err.Error()))
		}
		virtualImageRef = fmt.Sprintf("%s@%s", virtualImageRef, VirtualGoDigest.String())
	}

	VirtualImageRef, err = name.ParseReference(virtualImageRef)
	if err != nil {
		panic(fmt.Sprintf("%s/meta.VirtualImageRef is invalid: %s", Module, err.Error()))
	}
}

func VirtualImage() string {
	return VirtualImageRef.Name()
}
