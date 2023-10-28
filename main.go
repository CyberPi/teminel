package main

import (
	"fmt"

	"source.cyberpi.de/go/teminel/config"
)

func main() {
	path, err := config.SelectConfig()
	if err != nil {
		panic(err)
	}
	config, err := config.Load(path)
	if err != nil {
		panic(err)
	}
	for _, plugin := range config.Plugins {
		fmt.Println(plugin)
	}
}
