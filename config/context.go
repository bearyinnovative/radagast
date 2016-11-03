package config

import "context"

const KEY = "radagast:config"

func FromContext(c context.Context) map[string]interface{} {
	iconfig := c.Value(KEY)
	if config, ok := iconfig.(map[string]interface{}); ok {
		return config
	}

	panic("unable to get config from context")
}

func ToContext(c context.Context, config map[string]interface{}) context.Context {
	return context.WithValue(c, KEY, config)
}
