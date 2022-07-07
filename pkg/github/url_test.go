package github_test

import (
	"net/url"
	"testing"

	"github.com/frantjc/sequence/github"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	assert.Equal(t, "https://github.com", github.DefaultURL.String())
	assert.Equal(t, "https://api.github.com", github.DefaultAPIURL.String())
	assert.Equal(t, "https://api.github.com/graphql", github.DefaultGraphQLURL.String())
}

func TestGHESAPIURLFromBaseURL(t *testing.T) {
	var (
		base, err = url.Parse("https://github.myorg.com/")
		expected  = "https://github.myorg.com/api/v3"
	)
	assert.Nil(t, err)

	actual, err := github.APIURLFromBaseURL(base)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.String())
}

func TestGHESGraphQLURLFromBaseURL(t *testing.T) {
	var (
		base, err = url.Parse("https://github.myorg.com/")
		expected  = "https://github.myorg.com/api/v3/graphql"
	)
	assert.Nil(t, err)

	actual, err := github.GraphQLURLFromBaseURL(base)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.String())
}

func TestAPIURLFromBaseURL(t *testing.T) {
	var (
		base, err = url.Parse("https://github.com/")
		expected  = "https://api.github.com/"
	)
	assert.Nil(t, err)

	actual, err := github.APIURLFromBaseURL(base)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.String())
}

func TestGraphQLURLFromBaseURL(t *testing.T) {
	var (
		base, err = url.Parse("https://github.com/")
		expected  = "https://api.github.com/graphql"
	)
	assert.Nil(t, err)

	actual, err := github.GraphQLURLFromBaseURL(base)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.String())
}
