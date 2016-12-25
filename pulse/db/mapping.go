package db

import (
	"context"
	"errors"

	"gopkg.in/olivere/elastic.v5"
)

const mapping = `{
  "mappings": {
    "pullrequests": {
      "properties": {
        "id": {"type": "string", "index": "not_analyzed"},
        "number": {"type": "string", "index": "not_analyzed"},
        "url": {"type": "string", "index": "not_analyzed"},
        "state": {"type": "string", "index": "not_analyzed"},
        "title": {"type": "string"},
        "body": {"type": "string"},
        "created_at": {"type": "date"},
        "updated_at": {"type": "date"},
        "closed_at": {"type": "date"},
        "merged_at": {"type": "date"},
        "additions": {"type": "long"},
        "deletions": {"type": "long"},
        "changed_files": {"type": "long"},
        "repo": {
          "type": "nested",
          "properties": {
            "owner": {"type": "string", "index": "not_analyzed"},
            "name": {"type": "string", "index": "not_analyzed"}
          }
        },
        "user": {
          "type": "nested",
          "properties": {
            "login": {"type": "string", "index": "not_analyzed"},
            "name": {"type": "string", "index": "not_analyzed"}
          }
        },
        "merged_by": {
          "type": "nested",
          "properties": {
            "login": {"type": "string", "index": "not_analyzed"},
            "name": {"type": "string", "index": "not_analyzed"}
          }
        },
        "assignees": {
          "type": "nested",
          "properties": {
            "login": {"type": "string", "index": "not_analyzed"},
            "name": {"type": "string", "index": "not_analyzed"}
          }
        }
      }
    }
  }
}`

const PULSE_INDEX = "beary_pulse"

var errCreateFailed = errors.New("create mapping failed")

func CreateMapping(c context.Context, client *elastic.Client) error {
	rv, err := client.CreateIndex(PULSE_INDEX).BodyString(mapping).Do(c)
	if err != nil {
		return err
	}

	if !rv.Acknowledged {
		return errCreateFailed
	}

	return nil
}
