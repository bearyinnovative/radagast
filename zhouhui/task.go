package zhouhui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bearyinnovative/radagast/config"
	gh "github.com/bearyinnovative/radagast/github"
	"github.com/google/go-github/github"
)

const TaskName = "pullrequests"

var (
	ErrInvalidRepoSlug = errors.New("repo name should be `owner/name`")
	ErrDateRequired    = errors.New("date required, set with `ZHOUHUI_DATE`")
)

func ExecuteOnce(ctx context.Context) error {
	config := config.FromContext(ctx).Get("zhouhui").Config()
	repoSlug := strings.Split(config.Get("repo").String(), "/")
	if len(repoSlug) != 2 {
		return ErrInvalidRepoSlug
	}

	issue, err := buildIssue(config)
	if err != nil {
		return err
	}

	ghClient := gh.ClientFromContext(ctx)
	_, _, err = ghClient.Issues.Create(
		repoSlug[0],
		repoSlug[1],
		issue,
	)
	return err
}

func buildIssue(config config.Config) (*github.IssueRequest, error) {
	templateConfig := config.Get("issue-template").Config()

	date := os.Getenv("ZHOUHUI_DATE")
	if date == "" {
		return nil, ErrDateRequired
	}

	title := fmt.Sprintf(templateConfig.Get("title").String(), date)
	body := fmt.Sprintf(templateConfig.Get("body").String(), date)
	labels := strings.Split(templateConfig.Get("labels").String(), ",")

	return &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}, nil
}
