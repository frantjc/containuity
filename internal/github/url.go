package github

import "net/url"

var DefaultURL *url.URL

func init() {
	DefaultURL, _ = url.Parse("https://github.com/")
}
