package monitor_stale_issues

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/google/go-github/github"
	"github.com/hashicorp/go-multierror"
)

type repo struct {
	owner string
	name  string

	bearychatUserAliases map[string]string
	bearychatVchannelId  string

	ReportChan chan string
}

// TODO: better way to parse config
func getReposFromConfig(config map[string]interface{}) (repos []repo, err error) {
	irepos, ok := config["repos"].([]interface{})
	if !ok {
		err = errors.New("unable get repos")
		return
	}

	for _, irepo := range irepos {
		r, ok := irepo.(map[string]interface{})
		if !ok {
			err = errors.New("unable get repo")
			return
		}
		userAliases := make(map[string]string)
		for k, v := range r["bearychat-users"].(map[string]interface{}) {
			userAliases[k] = v.(string)
		}
		repoSlug := strings.Split(r["repo"].(string), "/")
		if len(repoSlug) != 2 {
			err = errors.New("repo name should be `owner/name`")
			return
		}
		repos = append(repos, repo{
			owner:                repoSlug[0],
			name:                 repoSlug[1],
			bearychatUserAliases: userAliases,
			bearychatVchannelId:  r["bearychat-vchannel-id"].(string),
			ReportChan:           make(chan string, 1024),
		})
	}

	return
}

func (r repo) String() string { return r.Slug() }
func (r repo) Slug() string   { return fmt.Sprintf("%s/%s", r.owner, r.name) }

func (r repo) SentReport(ctx context.Context) {
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
	for e := range checkErrChan {
		checkErr = multierror.Append(checkErr, e)
	}

	return checkErr.ErrorOrNil()
}

func checkRepo(c context.Context, github *github.Client, repo repo) error {
	if err := checkStalePullRequests(c, github, repo); err != nil {
		return err
	}

	return nil
}

type bearychatUser struct {
	name string
}

func getBearyChatUserFromGitHubUser(repo repo, ghUser *github.User) (u bearychatUser) {
	ghUserLogin := *ghUser.Login

	if name, present := repo.bearychatUserAliases[ghUserLogin]; present {
		u.name = name
	} else {
		u.name = ghUserLogin
	}

	return u
}
