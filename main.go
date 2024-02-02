package main

import (
	"os"

	"source.cyberpi.de/go/teminel/core/git"
	"source.cyberpi.de/go/teminel/core/proxy"
	"source.cyberpi.de/go/teminel/core/tmux"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "git":
			if err := git.Run(&git.Default, 8080); err != nil {
				panic(err)
			}
			return
		case "proxy":
			os.Args = os.Args[1:]
			proxy.Main()
			return
		}
	}
	err := tmux.Run()
	if err != nil {
		panic(err)
	}
}
