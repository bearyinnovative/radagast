package main

import (
	"context"
	"log"

	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/task"
	toml "github.com/pelletier/go-toml"
)

func main() {
	ctx := context.Background()
	ctx = getConfig(ctx)

	if err := task.Execute(ctx); err != nil {
		log.Fatalf("execute tasks failed: %+v", err)
	}
}

func getConfig(ctx context.Context) context.Context {
	c, err := toml.LoadFile("./radagast.toml")
	if err != nil {
		log.Fatalf("unable to get config: %+v", err)
		return nil
	}

	return config.ToContext(ctx, c.ToMap())
}
