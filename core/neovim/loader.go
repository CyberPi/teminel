package neovim

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/load"
)

type Loader struct {
	server *gitkit.Server
}

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

const workingDirectory = "/tmp/teminel/working"

func (cache *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if err := load.CloneBare("github.com", matches[2], bareDirectory, workingDirectory); err != nil {
			fmt.Println(err)
		}
	}
	cache.server.ServeHTTP(writer, request)
}
