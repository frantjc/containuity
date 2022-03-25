package github

import (
	"fmt"
	"net/url"
	"path"

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

	DefaultAPIURL, err = APIURLFromBaseURL(DefaultURL)
	if err != nil {
		panic(fmt.Sprintf("%s/github.DefaultAPIURL is not a valid URL", meta.Module))
	}

	DefaultGraphQLURL, err = GraphQLURLFromBaseURL(DefaultURL)
	if err != nil {
		panic(fmt.Sprintf("%s/github.DefaultGraphQLURL is not a valid URL", meta.Module))
	}
}

func APIURLFromBaseURL(base *url.URL) (*url.URL, error) {
	api, err := url.Parse(base.String())
	if err != nil {
		return nil, err
	}

	if api.Hostname() == "github.com" {
		api.Host = "api.github.com"
	} else {
		api.Path = path.Join(api.Path, "/api/v3")
	}
	return api, nil
}

func GraphQLURLFromBaseURL(base *url.URL) (*url.URL, error) {
	graphql, err := APIURLFromBaseURL(base)
	if err != nil {
		return nil, err
	}
	graphql.Path = path.Join(graphql.Path, "graphql")
	return graphql, nil
}
