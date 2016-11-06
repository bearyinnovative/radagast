package monitor_stale_issues

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bearyinnovative/radagast/config"
	gh "github.com/bearyinnovative/radagast/github"
	"github.com/google/go-github/github"
)

const TaskName = "monitor-stale-issues"

var (
	checkInterval = time.Duration(24 * time.Hour)
	//checkInterval = time.Duration(30 * time.Second)
)

type task struct {
	config       config.Config
	repos        []repo
	githubClient *github.Client
}

func makeTask(ctx context.Context) (*task, error) {
	config := config.FromContext(ctx).Get(TaskName).Config()

	repos, err := getReposFromConfig(config)
	if err != nil {
		return nil, err
	}

	github := gh.ClientFromContext(ctx)

	return &task{config, repos, github}, nil
}

func ExecuteOnce(ctx context.Context) error {
	task, err := makeTask(ctx)
	if err != nil {
		return err
	}

	return checkRepos(ctx, task.githubClient, task.repos)
}

func Execute(ctx context.Context) error {
	task, err := makeTask(ctx)
	if err != nil {
		return err
	}

	if err := checkRepos(ctx, task.githubClient, task.repos); err != nil {
		return err
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for range ticker.C {
		if err := checkRepos(ctx, task.githubClient, task.repos); err != nil {
			return err
		}
	}

	return nil
}

func logf(f string, args ...interface{}) {
	f = fmt.Sprintf("[%s] %s", TaskName, f)
	log.Printf(f, args...)
}
