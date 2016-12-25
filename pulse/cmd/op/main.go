package main

import (
	"context"
	"flag"
	"log"

	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/pulse/db"
)

var (
	reindex bool
)

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = db.MustMakeContext(ctx)

	dbClient := db.ClientFromContext(ctx)

	if reindex {
		if _, err := dbClient.DeleteIndex(db.PULSE_INDEX).Do(ctx); err != nil {
			log.Fatalf("delete index failed")
		} else {
			log.Printf("index deleted")
		}
	}

	if err := db.CreateMapping(ctx, dbClient); err == nil {
		log.Printf("created mapping")
	} else {
		log.Fatalf("create index failed: %+v", err)
	}
}

func init() {
	flag.BoolVar(&reindex, "reindex", false, "delete current index?")
}
