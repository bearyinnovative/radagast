package metric

import (
	"context"
	"errors"

	"github.com/bearyinnovative/radagast/pulse/db"

	"gopkg.in/olivere/elastic.v5"
)

var (
	errNotIndexed = errors.New("not indexed")
)

func IndexPullRequest(ctx context.Context, esClient *elastic.Client, pr *PullRequest) error {
	if pr == nil {
		return errNotIndexed
	}

	_, err := esClient.DeleteByQuery().
		Index(db.PULSE_INDEX).
		Type(TypePullRequest).
		Query(elastic.NewTermQuery("id", pr.ID)).
		Do(ctx)
	if err != nil {
		return err
	}

	_, err = esClient.Index().
		Index(db.PULSE_INDEX).
		Type(TypePullRequest).
		BodyJson(pr).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
