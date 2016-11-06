package pullrequests

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
	gh "github.com/bearyinnovative/radagast/github"
	"github.com/google/go-github/github"
)

var (
	ErrInvalidRepoSlug    = errors.New("repo name should be `owner/name`")
	ErrVchannelIdRequired = errors.New("`bearychat-vchannel-id` required")
)

// Repo represents a GitHub repoository.
type Repo struct {
	Owner string
	Name  string

	BCVchannelId string // target vchannel id

	ghClient      *github.Client
	bcUserAliases config.Config
}

// GetReposFromConfig build a list of repos from config.
func GetReposFromConfig(ctx context.Context, config config.Config) (repos []*Repo, err error) {
	ghClient := gh.ClientFromContext(ctx)
	bcUserAliases := bearychat.GitHubUsersFromContext(ctx)

	for _, r := range config.GetSlice("repos") {
		repoConfig := r.Config()
		repoSlug := strings.Split(repoConfig.Get("repo").String(), "/")
		if len(repoSlug) != 2 {
			err = ErrInvalidRepoSlug
			return
		}

		vchannelId := repoConfig.Get("bearychat-vchannel-id").String()
		if vchannelId == "" {
			err = ErrVchannelIdRequired
			return
		}

		repos = append(repos, &Repo{
			Owner: repoSlug[0],
			Name:  repoSlug[1],

			BCVchannelId: vchannelId,

			ghClient:      ghClient,
			bcUserAliases: bcUserAliases,
		})
	}

	return
}

func (r Repo) String() string { return r.Slug() }
func (r Repo) Slug() string   { return fmt.Sprintf("%s/%s", r.Owner, r.Name) }

func (r Repo) ghUserLoginToBCUserName(login string) string {
	s := r.bcUserAliases.Get(login).String()
	if s == "" {
		s = login
	}
	return fmt.Sprintf("@%s", s)
}
