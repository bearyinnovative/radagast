package main

import (
	"context"
	"log"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/github"
	task "github.com/bearyinnovative/radagast/monitor_stale_issues"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = bearychat.MustMakeContext(ctx)
	ctx = github.MustMakeContext(ctx)

	if err := task.ExecuteOnce(ctx); err != nil {
		log.Fatalf("execute tasks failed: %+v", err)
	}
}
