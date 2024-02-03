package git

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/utils"
)

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

func Main() {
	main()
}

func main() {
	host := utils.EnsureEnv("TEMINEL_HOST", "0.0.0.0:8080")
	flag.StringVar(&host, "host", host, "Server IP address to bind to.")

	flag.Parse()

	loader := Default

	err := utils.EnsureDirectories(loader.BareDirectory)
	if err != nil {
		panic(err)
	}
	mirrorConfig := gitkit.Config{
		Dir:        loader.BareDirectory,
		AutoCreate: true,
	}
	server := gitkit.New(mirrorConfig)
	if err := server.Setup(); err != nil {
		panic(err)
	}
	loader.server = server
	http.Handle("/", &loader)
	fmt.Println("Git mirror proxy started at:", host)
	panic(http.ListenAndServe(host, nil))
}
