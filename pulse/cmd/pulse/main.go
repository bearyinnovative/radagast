package main

import (
	"context"
	"log"

	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/github"
	"github.com/bearyinnovative/radagast/pulse/db"
	"github.com/bearyinnovative/radagast/pulse/metric"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = github.MustMakeContext(ctx)
	ctx = db.MustMakeContext(ctx)

	githubClient := github.ClientFromContext(ctx)
	prs, _, err := githubClient.PullRequests.List("bearyinnovative", "snitch", nil)
	if err != nil {
		panic(err)
	}

	dbClient := db.ClientFromContext(ctx)
	for _, pr := range prs {
		err := metric.IndexPullRequest(ctx, dbClient, metric.NewPullRequest(pr))
		if err != nil {
			log.Fatalf("index pr failed %+v", err)
		}
	}
}
