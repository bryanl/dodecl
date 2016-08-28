package util

import (
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// NewDigitalOceanClient creates an instance of godo.Client.
func NewDigitalOceanClient(pat string) *godo.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return godo.NewClient(tc)
}
