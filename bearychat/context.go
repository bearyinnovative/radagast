package bearychat

import (
	"context"

	bc "github.com/bcho/bearychat.go"
	"github.com/bearyinnovative/radagast/config"
)

const (
	RTM_CLIENT_KEY = "radagast:bearychat-rtm-client"
	CONFIG_KEY     = "bearychat-rtm-token"
)

func RTMClientFromContext(c context.Context) *bc.RTMClient {
	iclient := c.Value(RTM_CLIENT_KEY)
	if client, ok := iclient.(*bc.RTMClient); ok {
		return client
	}

	panic("unable to get bearychat rtm client from context")
}

func ToContext(c context.Context, client *bc.RTMClient) context.Context {
	return context.WithValue(c, RTM_CLIENT_KEY, client)
}

func MustMakeContext(ctx context.Context) context.Context {
	config := config.FromContext(ctx)
	rtmToken := config[CONFIG_KEY].(string)
	rtmClient, err := bc.NewRTMClient(rtmToken)
	if err != nil {
		panic(err)
	}

	return ToContext(ctx, rtmClient)
}
