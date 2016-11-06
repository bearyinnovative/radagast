package pullrequests

import (
	"context"
	"log"
	"time"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
)

const TaskName = "pullrequests"

type task struct {
	repos []*Repo
}

func newTask(ctx context.Context) (*task, error) {
	config := config.FromContext(ctx).Get("pullrequests").Config()
	repos, err := GetReposFromConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &task{repos}, nil
}

func ExecuteOnce(ctx context.Context) error {
	task, err := newTask(ctx)
	if err != nil {
		return err
	}

	for _, repo := range task.repos {
		log.Printf("checking %s", repo)

		stalePullRequets, err := repo.GetStalePullRequests(ctx)
		if err != nil {
			return err
		}
		for _, pullRequest := range stalePullRequets {
			report := pullRequest.Report()
			if pullRequest.IsStale() {
				log.Printf("%s is stale: %s", pullRequest, pullRequest.Type)
				sendReport(ctx, repo, report)
			} else {
				log.Printf("skipping %s: %s", pullRequest, report)
			}
		}
	}

	return nil
}

func sendReport(ctx context.Context, repo *Repo, report string) {
	bearychat.SendToVchannel(
		ctx,
		bearychat.RTMClientFromContext(ctx),
		bearychat.RTMMessage{
			Text:       report,
			VchannelId: repo.BCVchannelId,
			IsMarkdown: true,
		},
	)

	time.Sleep(1 * time.Second)
}
