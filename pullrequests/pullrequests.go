package pullrequests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

var (
	pullRequestFreshTime = time.Duration(24 * time.Hour)
	pullRequestStaleTime = time.Duration(72 * time.Hour)
)

type StaleType int

const (
	StaleTypeNormal StaleType = iota
	StaleTypeNoUpdates
	StaleTypeNoReviews
)

func (t StaleType) String() string {
	switch t {
	case StaleTypeNoReviews:
		return "no reviews"
	case StaleTypeNoUpdates:
		return "no updates"
	default:
		return "normal"
	}
}

// StalePullRequest represents a stale pr.
type StalePullRequest struct {
	Type        StaleType
	Repo        *Repo
	PullRequest *github.PullRequest
	reason      string
}

// CheckStalePullRequest checks pull request.
func CheckStalePullRequest(repo *Repo, pullRequest *github.PullRequest) (stalePullRequest *StalePullRequest, err error) {
	stalePullRequest = &StalePullRequest{
		Type:        StaleTypeNormal,
		Repo:        repo,
		PullRequest: pullRequest,
	}

	prTitle := *pullRequest.Title
	prState := *pullRequest.State

	if prState != "open" {
		stalePullRequest.reason = "pull request is not opened"
		return
	}

	if len(pullRequest.Assignees) < 1 {
		stalePullRequest.reason = "pull request has no assignees"
		return
	}

	if (*pullRequest.CreatedAt).Add(pullRequestFreshTime).After(time.Now()) {
		// GitHub's pr comment API seems like have a 24 hours lag,
		// so we have to check later.
		stalePullRequest.reason = "pull request is still fresh"
		return
	}

	if strings.HasPrefix(strings.ToLower(prTitle), "[wip]") {
		stalePullRequest.reason = "pull request is WIP"
		return
	}

	unreviewedAssignees, err := findUnreviewAssignees(repo, pullRequest)
	if err != nil {
		return
	}
	if len(unreviewedAssignees) > 0 {
		stalePullRequest.Type = StaleTypeNoReviews
		stalePullRequest.reason = strings.Join(unreviewedAssignees, " ")
		return
	}

	if (*pullRequest.UpdatedAt).Add(pullRequestStaleTime).Before(time.Now()) {
		stalePullRequest.Type = StaleTypeNoUpdates
		return
	}

	return
}

// IsStale indicates if the pull request is stale.
func (s StalePullRequest) IsStale() bool {
	switch s.Type {
	default:
		return true
	case StaleTypeNormal:
		return false
	}
}

// Report renders stale report for users.
func (s StalePullRequest) Report() string {
	switch s.Type {
	default:
		return s.reason
	case StaleTypeNoUpdates:
		pullRequestUser := *s.PullRequest.User
		return fmt.Sprintf(
			"%s 太久没有更新了 %s",
			s.reportPullRequest(),
			s.Repo.ghUserLoginToBCUserName(*pullRequestUser.Login),
		)
	case StaleTypeNoReviews:
		return fmt.Sprintf(
			"%s 还没有 review %s",
			s.reportPullRequest(),
			s.reason,
		)
	}
}

func (s StalePullRequest) reportPullRequest() string {
	return fmt.Sprintf(
		"`%s` PR  [#%d](https://github.com/%s/pull/%d): %s\n\n",
		s.Repo.Name,
		*s.PullRequest.Number,
		s.Repo.Slug(),
		*s.PullRequest.Number,
		*s.PullRequest.Title,
	)
}

func (s StalePullRequest) String() string {
	return fmt.Sprintf(
		"%s PR [#%d]",
		s.Repo.Name,
		*s.PullRequest.Number,
	)
}

var stalePullRequestListOpts = &github.PullRequestListOptions{
	State:     "open",
	Direction: "desc",
}

// GetStalePullRequests returns all stale pull requests for this repo.
func (r *Repo) GetStalePullRequests(ctx context.Context) ([]*StalePullRequest, error) {
	pullRequests, _, err := r.ghClient.PullRequests.List(
		r.Owner,
		r.Name,
		stalePullRequestListOpts,
	)
	if err != nil {
		return nil, err
	}

	var stalePullRequests []*StalePullRequest
	for _, pullRequest := range pullRequests {
		stalePullRequest, err := CheckStalePullRequest(r, pullRequest)
		if err != nil {
			return nil, err
		}
		stalePullRequests = append(stalePullRequests, stalePullRequest)
	}

	return stalePullRequests, nil
}

func findUnreviewAssignees(repo *Repo, pullRequest *github.PullRequest) (bcUserNames []string, err error) {
	comments, _, err := repo.ghClient.PullRequests.ListComments(
		repo.Owner,
		repo.Name,
		*pullRequest.Number,
		nil,
	)
	if err != nil {
		return nil, err
	}

	unreviewAssignees := make(map[int]string)
	for _, assignee := range pullRequest.Assignees {
		unreviewAssignees[*assignee.ID] = *assignee.Login
	}
	for _, comment := range comments {
		commentUser := *comment.User
		delete(unreviewAssignees, *commentUser.ID)
	}

	for _, login := range unreviewAssignees {
		bcUserName := repo.ghUserLoginToBCUserName(login)
		bcUserNames = append(bcUserNames, bcUserName)
	}

	return
}
