package main

import (
	"fmt"
	"os"

	"source.cyberpi.de/go/teminel/core/git"
	"source.cyberpi.de/go/teminel/core/proxy"
	"source.cyberpi.de/go/teminel/core/tmux"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "git":
			fmt.Println("Starting teminel in git mirror mode")
			os.Args = os.Args[1:]
			git.Main()
			return
		case "proxy":
			fmt.Println("Starting teminel in http proxy mode")
			os.Args = os.Args[1:]
			proxy.Main()
			return
		}
	}
	fmt.Println("Starting teminel in tmux plugin manager mode")
	tmux.Main()
}
