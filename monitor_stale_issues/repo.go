package monitor_stale_issues

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
	"github.com/google/go-github/github"
	"github.com/hashicorp/go-multierror"
)

type repo struct {
	owner string
	name  string

	bearychatVchannelId string

	ReportChan chan string
}

func getReposFromConfig(config config.Config) (repos []repo, err error) {
	for _, irepo := range config.GetSlice("repos") {
		repoConfig := irepo.Config()
		repoSlug := strings.Split(repoConfig.Get("repo").String(), "/")
		if len(repoSlug) != 2 {
			err = errors.New("repo name should be `owner/name`")
			return
		}
		repos = append(repos, repo{
			owner:               repoSlug[0],
			name:                repoSlug[1],
			bearychatVchannelId: repoConfig.Get("bearychat-vchannel-id").String(),
			ReportChan:          make(chan string, 1024),
		})
	}

	return
}

func (r repo) String() string { return r.Slug() }
func (r repo) Slug() string   { return fmt.Sprintf("%s/%s", r.owner, r.name) }

func (r repo) SendReport(ctx context.Context) {
	bc := bearychat.RTMClientFromContext(ctx)
	for report := range r.ReportChan {
		bearychat.SendToVchannel(
			ctx,
			bc,
			bearychat.RTMMessage{
				Text:       report,
				VchannelId: r.bearychatVchannelId,
				IsMarkdown: true,
			},
		)
		time.Sleep(1 * time.Second)
	}
}

func checkRepos(ctx context.Context, github *github.Client, repos []repo) error {
	checkErrChan := make(chan error, len(repos))
	for _, r := range repos {
		go func(r repo) {
			checkErrChan <- checkRepo(ctx, github, r)
		}(r)
	}

	var checkErr *multierror.Error
	for i := 0; i < len(repos); i++ {
		if err := <-checkErrChan; err != nil {
			checkErr = multierror.Append(checkErr, err)
		}
	}

	return checkErr.ErrorOrNil()
}

func checkRepo(c context.Context, github *github.Client, repo repo) error {
	go repo.SendReport(c)

	if err := checkStalePullRequests(c, github, repo); err != nil {
		return err
	}

	return nil
}

type bearychatUser struct {
	name string
}

func getBearyChatUserFromGitHubUser(c context.Context, ghUser *github.User) (u bearychatUser) {
	ghUsers := bearychat.GitHubUsersFromContext(c)
	ghUserName := *ghUser.Login
	u.name = ghUsers.Get(ghUserName).String()
	if u.name == "" {
		u.name = ghUserName
	}

	return
}
