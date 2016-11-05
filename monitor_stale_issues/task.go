package monitor_stale_issues

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/bearyinnovative/radagast/config"
	"github.com/google/go-github/github"
)

const TaskName = "monitor-stale-issues"

var (
	checkInterval = time.Duration(24 * time.Hour)
	//checkInterval = time.Duration(30 * time.Second)

	checkableDuration = time.Duration(24 * time.Hour)
)

type repo struct {
	owner string
	name  string

	bearychatUserAliases map[string]string
	bearychatVchannelId  string
}

func (r repo) String() string { return fmt.Sprintf("%s/%s", r.owner, r.name) }

func Execute(ctx context.Context) error {
	config := getConfig(ctx)

	repos, err := getRepos(config)
	if err != nil {
		return err
	}

	github, err := getGitHubClient(config)
	if err != nil {
		return err
	}

	checkRepos(github, repos)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		checkRepos(github, repos)
	}

	return nil
}

func checkRepos(github *github.Client, repos []repo) error {
	for _, repo := range repos {
		if err := checkRepo(github, repo); err != nil {
			return err
		}
	}

	return nil
}

func checkRepo(github *github.Client, repo repo) error {
	if err := checkStalePullRequests(github, repo); err != nil {
		return err
	}

	return nil
}

var stalePullRequestListOpts = &github.PullRequestListOptions{
	State:     "open",
	Direction: "desc",
}

func checkStalePullRequests(github *github.Client, repo repo) error {
	logf("checking stale pr for %s", repo)

	prs, _, err := github.PullRequests.List(
		repo.owner,
		repo.name,
		stalePullRequestListOpts,
	)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		if err := checkStalePullRequest(github, repo, pr); err != nil {
			return err
		}
	}

	return nil
}

func checkStalePullRequest(githubClient *github.Client, repo repo, pr *github.PullRequest) (err error) {
	prNumber := *pr.Number
	prTitle := *pr.Title
	prState := *pr.State

	// fast path
	if prState != "open" {
		return
	}
	if (pr.Mergeable != nil && !*pr.Mergeable) ||
		(pr.Merged != nil && *pr.Merged) {
		return
	}
	if len(pr.Assignees) < 1 {
		return
	}
	if (*pr.CreatedAt).Add(checkableDuration).After(time.Now()) {
		// GitHub's pr comment API seems like have a 24 hours lag,
		// so we have to check later.
		return
	}
	if strings.HasPrefix(strings.ToLower(prTitle), "[wip]") {
		return
	}

	logf("checking pull request: [%d][%s] %s", prNumber, prState, prTitle)

	// TODO: cache stats
	comments, _, err := githubClient.PullRequests.ListComments(
		repo.owner,
		repo.name,
		prNumber,
		nil,
	)
	if err != nil {
		return
	}

	unreviewedAssignees := make(map[int]*github.User)
	for _, assignee := range pr.Assignees {
		unreviewedAssignees[*assignee.ID] = assignee
	}
	for _, comment := range comments {
		commentUser := *comment.User
		delete(unreviewedAssignees, *commentUser.ID)
	}

	for _, unreviewedAssignee := range unreviewedAssignees {
		logf("unreviewed assignee: %s", *unreviewedAssignee.Login)
	}

	return nil
}

func getConfig(ctx context.Context) map[string]interface{} {
	config := config.FromContext(ctx)
	return config[TaskName].(map[string]interface{})
}

// TODO: better way to parse config
func getRepos(config map[string]interface{}) (repos []repo, err error) {
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
