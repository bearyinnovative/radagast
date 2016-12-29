package airtable

import (
	"context"

	"github.com/bearyinnovative/radagast/config"
	"github.com/fabioberger/airtable-go"
)

const (
	KEY_AIRTABLE_CLIENT = "radagast:airtable.client"
)

func ClientFromContext(c context.Context) *airtable.Client {
	iclient := c.Value(KEY_AIRTABLE_CLIENT)
	if client, ok := iclient.(*airtable.Client); ok {
		return client
	}

	panic("unable to get airtable client from context")
}

func MustMakeContext(c context.Context) context.Context {
	airtableConfig := config.FromContext(c).Get("airtable").Config()
	apiKey := airtableConfig.Get("api-key").String()
	base := airtableConfig.Get("base").String()
	client := airtable.New(apiKey, base, false)

	return context.WithValue(c, KEY_AIRTABLE_CLIENT, client)
}
