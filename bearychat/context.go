package bearychat

import (
	"context"

	bc "github.com/bcho/bearychat.go"
	"github.com/bearyinnovative/radagast/config"
)

const (
	KEY_RTM_CLIENT   = "radagast:bearychat.rtm-client"
	KEY_GITHUB_USERS = "radagast:bearychat.github-users"
)

func RTMClientFromContext(c context.Context) *bc.RTMClient {
	iclient := c.Value(KEY_RTM_CLIENT)
	if client, ok := iclient.(*bc.RTMClient); ok {
		return client
	}

	panic("unable to get bearychat rtm client from context")
}

func GitHubUsersFromContext(c context.Context) config.Config {
	iusers := c.Value(KEY_GITHUB_USERS)
	if users, ok := iusers.(config.Config); ok {
		return users
	}

	panic("unable to get bearychat github users from context")
}

func MustMakeContext(c context.Context) context.Context {
	config := config.FromContext(c)

	rtmToken := config.Get("bearychat.rtm-token").String()
	rtmClient, err := bc.NewRTMClient(rtmToken)
	if err != nil {
		panic(err)
	}
	c = context.WithValue(c, KEY_RTM_CLIENT, rtmClient)

	users := config.Get("bearychat.github-users").Config()
	c = context.WithValue(c, KEY_GITHUB_USERS, users)

	return c
}
