package monitor_stale_issues

import (
	"strings"
	"time"

	"github.com/google/go-github/github"
)

var checkableDuration = time.Duration(24 * time.Hour)

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

func checkStalePullRequest(githubClient *github.Client, repo repo, pr *github.PullRequest) (err error) {
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

	for _, unreviewedAssignee := range unreviewedAssignees {
		logf("unreviewed assignee: %s", *unreviewedAssignee.Login)
	}

	return nil
}
