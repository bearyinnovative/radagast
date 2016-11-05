package main

import (
	"context"
	"log"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/task"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = bearychat.MustMakeContext(ctx)

	if err := task.Execute(ctx); err != nil {
		log.Fatalf("execute tasks failed: %+v", err)
	}
}
