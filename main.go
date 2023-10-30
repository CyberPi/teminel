package main

import (
	"source.cyberpi.de/go/teminel/core/tmux"
)

func main() {
	configTmux := &tmux.Config{}
	path, err := tmux.SelectConfig()
	if err != nil {
		panic(err)
	}
	err = configTmux.Read(path)
	if err != nil {
		panic(err)
	}
	configTmux.Install()
	configTmux.Load()
}
