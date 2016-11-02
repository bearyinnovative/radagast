package main

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

func main() {
	config, err := toml.LoadFile("./radagast.toml")
	if err != nil {
		panic(err)
	}

	tasks := config.Get("tasks").([]interface{})
	fmt.Printf("%+v\n", tasks)

	repos := config.Get("monitor-stale-issues.repos").([]*toml.TomlTree)
	fmt.Printf("%+v\n", repos)
}
