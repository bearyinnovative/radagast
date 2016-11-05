package monitor_stale_issues

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-multierror"
)

type repo struct {
	owner string
	name  string

	bearychatUserAliases map[string]string
	bearychatVchannelId  string
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
		})
	}

	return
}

func (r repo) String() string { return fmt.Sprintf("%s/%s", r.owner, r.name) }

func checkRepos(github *github.Client, repos []repo) error {
	checkErrChan := make(chan error, len(repos))
	for _, r := range repos {
		go func(r repo) {
			checkErrChan <- checkRepo(github, r)
		}(r)
	}

	var checkErr *multierror.Error
	for e := range checkErrChan {
		checkErr = multierror.Append(checkErr, e)
	}

	return checkErr.ErrorOrNil()
}

func checkRepo(github *github.Client, repo repo) error {
	if err := checkStalePullRequests(github, repo); err != nil {
		return err
	}

	return nil
}
