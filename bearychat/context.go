package bearychat

import (
	"context"

	bc "github.com/bcho/bearychat.go"
	"github.com/bearyinnovative/radagast/config"
)

const (
	KEY_RTM_CLIENT = "radagast:bearychat.rtm-client"
	KEY_USERS      = "radagast:bearychat.users"
)

func RTMClientFromContext(c context.Context) *bc.RTMClient {
	iclient := c.Value(KEY_RTM_CLIENT)
	if client, ok := iclient.(*bc.RTMClient); ok {
		return client
	}

	panic("unable to get bearychat rtm client from context")
}

func UsersFromContext(c context.Context) config.Config {
	iusers := c.Value(KEY_USERS)
	if users, ok := iusers.(config.Config); ok {
		return users
	}

	panic("unable to get bearychat users from context")
}

func MustMakeContext(c context.Context) context.Context {
	config := config.FromContext(c)

	rtmToken := config.Get("bearychat.rtm-token").String()
	rtmClient, err := bc.NewRTMClient(rtmToken)
	if err != nil {
		panic(err)
	}
	c = context.WithValue(c, KEY_RTM_CLIENT, rtmClient)

	users := config.Get("bearychat.users").Config()
	c = context.WithValue(c, KEY_USERS, users)

	return c
}
