package main

import (
	"context"
	"log"

	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/pulse/db"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = db.MustMakeContext(ctx)

	dbClient := db.ClientFromContext(ctx)
	if err := db.CreateMapping(ctx, dbClient); err == nil {
		log.Printf("created mapping")
	} else {
		log.Fatalf("create index failed: %+v", err)
	}
}
