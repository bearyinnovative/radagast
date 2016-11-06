package config

import (
	"context"

	toml "github.com/pelletier/go-toml"
)

const KEY = "radagast:config"

func FromContext(c context.Context) Config {
	iconfig := c.Value(KEY)
	if config, ok := iconfig.(Config); ok {
		return config
	}

	panic("unable to get config from context")
}

func ToContext(c context.Context, config Config) context.Context {
	return context.WithValue(c, KEY, config)
}

func MustMakeContext(ctx context.Context, configPath string) context.Context {
	c, err := toml.LoadFile(configPath)
	if err != nil {
		panic(err)
	}

	return ToContext(ctx, NewFromMap(c.ToMap()))
}
