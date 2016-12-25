package metric

import (
	"time"

	"github.com/google/go-github/github"
)

type User struct {
	Login *string `json:"login"`
	Name  *string `json:"name"`
}

func NewUser(user *github.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		Login: user.Login,
		Name:  user.Name,
	}
}

type Repo struct {
	Owner *string `json:"owner"`
	Name  *string `json:"name"`
}

func NewRepo(repo *github.Repository) *Repo {
	if repo == nil {
		return nil
	}

	return &Repo{
		Owner: repo.Owner.Login,
		Name:  repo.Name,
	}
}

type PullRequest struct {
	ID           *int       `json:"id"`
	Number       *int       `json:"number"`
	URL          *string    `json:"url"`
	State        *string    `json:"state"`
	Title        *string    `json:"title"`
	Body         *string    `json:"body"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	MergedAt     *time.Time `json:"merged_at"`
	Additions    *int       `json:"additions"`
	Deletions    *int       `json:"deletions"`
	ChangedFiles *int       `json:"changed_files"`

	Repo      *Repo   `json:"repo"`
	User      *User   `json:"user"`
	MergedBy  *User   `json:"merged_by"`
	Assignees []*User `json:"assignees"`
}

func NewPullRequest(pr *github.PullRequest) *PullRequest {
	if pr == nil {
		return nil
	}

	pullRequest := &PullRequest{
		ID:           pr.ID,
		Number:       pr.Number,
		URL:          pr.URL,
		State:        pr.State,
		Title:        pr.Title,
		Body:         pr.Body,
		CreatedAt:    pr.CreatedAt,
		UpdatedAt:    pr.UpdatedAt,
		ClosedAt:     pr.ClosedAt,
		MergedAt:     pr.MergedAt,
		Additions:    pr.Additions,
		Deletions:    pr.Deletions,
		ChangedFiles: pr.ChangedFiles,

		User:     NewUser(pr.User),
		MergedBy: NewUser(pr.MergedBy),
	}

	if pr.Base != nil {
		pullRequest.Repo = NewRepo(pr.Base.Repo)
	}

	for _, a := range pr.Assignees {
		pullRequest.Assignees = append(pullRequest.Assignees, NewUser(a))
	}

	return pullRequest
}
