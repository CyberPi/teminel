package main

import (
	"os"

	"source.cyberpi.de/go/teminel/core/neovim"
	"source.cyberpi.de/go/teminel/core/tmux"
)

func main() {
	if len(os.Args) == 1 {
		err := tmux.Run()
		if err != nil {
			panic(err)
		}
	}
	err := neovim.Run()
	if err != nil {
		panic(err)
	}
}
