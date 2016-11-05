package monitor_stale_issues

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-multierror"
)

var checkableDuration = time.Duration(24 * time.Hour)

var stalePullRequestListOpts = &github.PullRequestListOptions{
	State:     "open",
	Direction: "desc",
}

func checkStalePullRequests(c context.Context, githubClient *github.Client, repo repo) error {
	logf("checking stale pr for %s", repo)

	prs, _, err := githubClient.PullRequests.List(
		repo.owner,
		repo.name,
		stalePullRequestListOpts,
	)
	if err != nil {
		return err
	}

	checkErrChan := make(chan error, len(prs))
	for _, pr := range prs {
		go func(pr *github.PullRequest) {
			checkErrChan <- checkStalePullRequest(c, githubClient, repo, pr)
		}(pr)
	}

	var checkErr *multierror.Error
	for e := range checkErrChan {
		checkErr = multierror.Append(checkErr, e)
	}

	return checkErr.ErrorOrNil()
}

func isStalePullRequest(pr *github.PullRequest) (bool, string) {
	prTitle := *pr.Title
	prState := *pr.State

	if prState != "open" {
		return false, "pr is not opened"
	}
	if len(pr.Assignees) < 1 {
		return false, "no assignees found"
	}
	if (*pr.CreatedAt).Add(checkableDuration).After(time.Now()) {
		// GitHub's pr comment API seems like have a 24 hours lag,
		// so we have to check later.
		return false, "still a fresh pr, check later"
	}
	if strings.HasPrefix(strings.ToLower(prTitle), "[wip]") {
		return false, "still a WIP pr"
	}

	return true, ""
}

func checkStalePullRequest(c context.Context, githubClient *github.Client, repo repo, pr *github.PullRequest) (err error) {
	prNumber := *pr.Number
	prTitle := *pr.Title

	if isStale, reason := isStalePullRequest(pr); !isStale {
		logf(
			"skipping pull request [%d] %s: %s",
			prNumber,
			prTitle,
			reason,
		)
		return nil
	}

	logf("checking pull request: [%d] %s", prNumber, prTitle)

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

	var unreviewedUsers []bearychatUser
	for _, assignee := range unreviewedAssignees {
		unreviewedUsers = append(unreviewedUsers, getBearyChatUserFromGitHubUser(repo, assignee))
	}

	return notifyUnreviewdAssignees(c, repo, pr, unreviewedUsers)
}

func notifyUnreviewdAssignees(c context.Context, repo repo, pr *github.PullRequest, users []bearychatUser) error {
	if len(users) < 1 {
		return nil
	}

	var report string

	prNumber := *pr.Number
	prTitle := *pr.Title
	report = fmt.Sprintf(
		"`%s` PR [#%d](https://github.com/%s/pull/%d):%s\n \n还没有 review",
		repo.name,
		prNumber,
		repo.Slug(),
		prNumber,
		prTitle,
	)

	for _, user := range users {
		report = fmt.Sprintf("%s @%s", report, user.name)
	}

	repo.ReportChan <- report

	return nil
}
