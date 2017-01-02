package worker

import (
	"context"
	"log"
	"time"

	"github.com/bearyinnovative/radagast/pulse/metric"
	"github.com/google/go-github/github"
	"gopkg.in/olivere/elastic.v5"
)

var (
	pullRequestIndexerInterval = 15 * time.Minute
)

// TODO cancel channel
func RunPullRequestIndexer(repo metric.Repo, githubClient *github.Client, esClient *elastic.Client) error {
	ticker := time.NewTicker(pullRequestIndexerInterval)
	defer ticker.Stop()

	ctx := context.Background()
	if err := indexPullRequests(ctx, repo, githubClient, esClient); err != nil {
		return err
	}

	for range ticker.C {
		ctx := context.Background()
		if err := indexPullRequests(ctx, repo, githubClient, esClient); err != nil {
			return err
		}
	}

	return nil
}

const (
	pagesToList = 3
	perPage     = 25
)

// TODO context
func indexPullRequests(ctx context.Context, repo metric.Repo, githubClient *github.Client, esClient *elastic.Client) error {
	log.Printf("start indexing pull requests for %s", repo)

	for page := 1; page < pagesToList; page++ {
		log.Printf("listing page %d for %s", page, repo)

		prs, err := listPage(page, repo, githubClient)
		if err != nil {
			return nil
		}

		for _, pr := range prs {
			pullRequest := metric.NewPullRequest(pr)
			err := metric.IndexPullRequest(ctx, esClient, pullRequest)
			if err != nil {
				return err
			}
			log.Printf("pull request %s indexed", pullRequest)
		}

		time.Sleep(3 * time.Second)
	}

	log.Printf("pull requests indexing for %s finished", repo)
	return nil
}

func listPage(page int, repo metric.Repo, githubClient *github.Client) ([]*github.PullRequest, error) {
	listOpts := &github.PullRequestListOptions{
		State: "all",
		Sort:  "updated",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	prs, _, err := githubClient.PullRequests.List(*repo.Name, *repo.Name, listOpts)
	return prs, err
}
