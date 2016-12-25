package metric

import "time"

type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type PullRequest struct {
	ID           string    `json:"id"`
	Number       string    `json:"number"`
	URL          string    `json:"url"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	MergedAt     time.Time `json:"merged_at"`
	Additions    int       `json:"additions"`
	Deletions    int       `json:"deletions"`
	ChangedFiles int       `json:"changed_files"`

	Repo      *Repo   `json:"repo"`
	User      *User   `json:"user"`
	MergedBy  *User   `json:"merged_by"`
	Assignees []*User `json:"assignees"`
}
