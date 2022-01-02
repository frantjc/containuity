package github

import (
	"fmt"
	"net/url"

	"github.com/frantjc/sequence"
)

var (
	URL        *url.URL
	APIURL     *url.URL
	GraphQLURL *url.URL
)

func init() {
	var err error
	URL, err = url.Parse("https://github.com")
	if err != nil {
		panic(fmt.Sprintf("%s/internal/github.URL is not a valid URL", sequence.Module))
	}

	APIURL, err = url.Parse("https://api.github.com")
	if err != nil {
		panic(fmt.Sprintf("%s/internal/github.APIURL is not a valid URL", sequence.Module))
	}

	GraphQLURL, err = url.Parse("https://api.github.com/graphql")
	if err != nil {
		panic(fmt.Sprintf("%s/internal/github.GraphQLURL is not a valid URL", sequence.Module))
	}
}
