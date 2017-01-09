package db

import (
	"context"

	"github.com/bearyinnovative/radagast/config"

	"gopkg.in/olivere/elastic.v5"
)

const CLIENT_KEY = "db_client"

func ClientFromContext(c context.Context) *elastic.Client {
	client, ok := c.Value(CLIENT_KEY).(*elastic.Client)

	if !ok {
		panic("db.client required")
	}

	return client
}

func ClientToContext(c context.Context, client *elastic.Client) context.Context {
	return context.WithValue(c, CLIENT_KEY, client)
}

func MustMakeContext(c context.Context) context.Context {
	config := config.FromContext(c).Get("pulse.db").Config()
	dbUrl := config.Get("url").String()
	if dbUrl == "" {
		panic("`pulse.db` required")
	}

	client, err := elastic.NewClient(
		elastic.SetURL(config.Get("url").String()),
	)
	if err != nil {
		panic(err)
	}

	return ClientToContext(c, client)
}
