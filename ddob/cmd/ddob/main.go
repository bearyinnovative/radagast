package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/ddob"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = bearychat.MustMakeContext(ctx)

	config := config.FromContext(ctx)

	p := regexp.MustCompile(config.Get("ddob.pattern").String())
	addrs, err := ddob.ListInterfaceAddrs(p)
	if err != nil {
		log.Fatal(err)
	}

	if len(addrs) < 1 {
		log.Fatalf("no addrs returned")
	}

	bearychat.SendToVchannel(
		ctx,
		bearychat.RTMClientFromContext(ctx),
		bearychat.RTMMessage{
			Text:       fmt.Sprintf("`ddob` 我的 IP 是 `%s`", addrs[0]),
			VchannelId: config.Get("ddob.bearychat-vchannel-id").String(),
		},
	)
}
