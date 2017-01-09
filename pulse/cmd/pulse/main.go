package main

import (
	"context"
	"log"

	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/github"
	"github.com/bearyinnovative/radagast/pulse/app"
	"github.com/bearyinnovative/radagast/pulse/db"
	"github.com/bearyinnovative/radagast/pulse/metric"
	"github.com/bearyinnovative/radagast/pulse/worker"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = github.MustMakeContext(ctx)
	ctx = db.MustMakeContext(ctx)

	config := config.FromContext(ctx).Get("pulse").Config()
	for _, r := range config.GetSlice("repos") {
		repoConfig := r.Config()
		repo := metric.NewRepoFromString(
			repoConfig.Get("owner").String(),
			repoConfig.Get("name").String(),
		)
		go indexRepo(ctx, *repo)

	}

	app.Serve(ctx)
}

func indexRepo(ctx context.Context, repo metric.Repo) {
	githubClient := github.ClientFromContext(ctx)
	esClient := db.ClientFromContext(ctx)

	log.Fatal(worker.RunPullRequestIndexer(repo, githubClient, esClient))
}
