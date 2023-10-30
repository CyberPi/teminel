package main

import (
	"source.cyberpi.de/go/teminel/core/tmux"
)

func main() {
	err := tmux.Run()
	if err != nil {
		panic(err)
	}
}
