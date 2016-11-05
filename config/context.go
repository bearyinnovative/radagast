package config

import (
	"context"

	toml "github.com/pelletier/go-toml"
)

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

func MustMakeContext(ctx context.Context, configPath string) context.Context {
	c, err := toml.LoadFile(configPath)
	if err != nil {
		panic(err)
	}

	return ToContext(ctx, c.ToMap())
}
