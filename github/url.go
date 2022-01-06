package github

import (
	"fmt"
	"net/url"

	"github.com/frantjc/sequence/meta"
)

var (
	DefaultURL        *url.URL
	DefaultAPIURL     *url.URL
	DefaultGraphQLURL *url.URL
)

func init() {
	var err error
	DefaultURL, err = url.Parse("https://github.com")
	if err != nil {
		panic(fmt.Sprintf("%s/github.DefaultURL is not a valid URL", meta.Module))
	}

	DefaultAPIURL, err = url.Parse("https://api.github.com")
	if err != nil {
		panic(fmt.Sprintf("%s/github.DefaultAPIURL is not a valid URL", meta.Module))
	}

	DefaultGraphQLURL, err = url.Parse("https://api.github.com/graphql")
	if err != nil {
		panic(fmt.Sprintf("%s/github.DefaultGraphQLURL is not a valid URL", meta.Module))
	}
}
