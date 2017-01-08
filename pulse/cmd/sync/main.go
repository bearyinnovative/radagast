package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/bearyinnovative/radagast/config"
	"github.com/bearyinnovative/radagast/github"
	"github.com/bearyinnovative/radagast/pulse/db"
	"github.com/bearyinnovative/radagast/pulse/metric"
	gogithub "github.com/google/go-github/github"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = github.MustMakeContext(ctx)
	ctx = db.MustMakeContext(ctx)

	config := config.FromContext(ctx).Get("pulse").Config()

	var wg sync.WaitGroup

	for _, r := range config.GetSlice("repos") {
		repoConfig := r.Config()
		repo := metric.NewRepoFromString(
			repoConfig.Get("owner").String(),
			repoConfig.Get("name").String(),
		)
		wg.Add(1)
		go indexRepo(ctx, *repo, &wg)

	}

	wg.Wait()
	log.Printf("all repo synced")
}

func indexRepo(ctx context.Context, repo metric.Repo, wg *sync.WaitGroup) {
	defer wg.Done()

	githubClient := github.ClientFromContext(ctx)
	esClient := db.ClientFromContext(ctx)

	listOpts := &gogithub.PullRequestListOptions{
		State:     "all",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: gogithub.ListOptions{
			PerPage: 25,
		},
	}

	for {
		prs, resp, err := githubClient.PullRequests.List(*repo.Owner, *repo.Name, listOpts)
		if err != nil {
			log.Printf("%s sync failed: %+v", repo, err)
			return
		}

		for _, pr := range prs {
			pullRequest := metric.NewPullRequest(pr)
			err := metric.IndexPullRequest(ctx, esClient, pullRequest)
			if err != nil {
				log.Printf("%s sync failed: %+v", repo, err)
			}
			log.Printf("pull request %s indexed", pullRequest)
		}

		if resp.NextPage == 0 {
			log.Printf("%s sync finished", repo)
			return
		}

		listOpts.ListOptions.Page = resp.NextPage
		time.Sleep(3 * time.Second)
	}
}
