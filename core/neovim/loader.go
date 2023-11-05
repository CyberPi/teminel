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

var repositoryMatcher = regexp.MustCompile(`((^.+?)\.git).*$`)

func (cache *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if err := load.CloneBare("https://hasdlkfhasdf.com", mirrorDir, matches[2]); err != nil {
			fmt.Println(err)
		}
	}
	cache.server.ServeHTTP(writer, request)
}
