package monitor_stale_issues

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"

	"github.com/bearyinnovative/radagast/config"
	"github.com/google/go-github/github"
)

const TaskName = "monitor-stale-issues"

var (
	checkInterval = time.Duration(24 * time.Hour)
	//checkInterval = time.Duration(30 * time.Second)
)

func Execute(ctx context.Context) error {
	config := getConfig(ctx)

	repos, err := getReposFromConfig(config)
	if err != nil {
		return err
	}

	github, err := getGitHubClient(config)
	if err != nil {
		return err
	}

	if err := checkRepos(github, repos); err != nil {
		return err
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for range ticker.C {
		if err := checkRepos(github, repos); err != nil {
			return err
		}
	}

	return nil
}

func getConfig(ctx context.Context) map[string]interface{} {
	config := config.FromContext(ctx)
	return config[TaskName].(map[string]interface{})
}

func getGitHubClient(config map[string]interface{}) (*github.Client, error) {
	token, ok := config["github-token"].(string)
	if !ok || token == "" {
		return nil, errors.New("`github-token` is required")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc), nil
}

func logf(f string, args ...interface{}) {
	f = fmt.Sprintf("[%s] %s", TaskName, f)
	log.Printf(f, args...)
}
