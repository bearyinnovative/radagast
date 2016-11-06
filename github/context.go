package github

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/bearyinnovative/radagast/config"
	"github.com/google/go-github/github"
)

const (
	KEY_GITHUB_CLIENT = "radagast:github.client"
)

func ClientFromContext(c context.Context) *github.Client {
	iclient := c.Value(KEY_GITHUB_CLIENT)
	if client, ok := iclient.(*github.Client); ok {
		return client
	}

	panic("unable to get github client from context")
}

func MustMakeContext(c context.Context) context.Context {
	token := config.FromContext(c).Get("github.api-token").String()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return context.WithValue(c, KEY_GITHUB_CLIENT, client)
}
