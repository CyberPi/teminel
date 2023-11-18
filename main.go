package main

import (
	"flag"

	"source.cyberpi.de/go/teminel/core/neovim"
	"source.cyberpi.de/go/teminel/core/tmux"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 1 {
		switch flag.Args()[0] {
		case "git":
			if err := neovim.Run(); err != nil {
				panic(err)
			}
			return

		case "proxy":
			if err := proxy.Run(); err != nil {
				panic(err)
			}
			return
		}
	}
	err := tmux.Run()
	if err != nil {
		panic(err)
	}
}
