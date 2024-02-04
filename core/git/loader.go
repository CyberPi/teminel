package git

import (
	"fmt"
	"net/http"
	"path"
	"regexp"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/load"
)

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

type Loader struct {
	Source           *load.GitSource
	BareDirectory    string
	WorkingDirectory string
	HomeDirectory    string
	server           *gitkit.Server
}

func (repository *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("New request:", request.URL.Path, "Method:", request.Method)
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if len(matches) >= 3 {
			if err := repository.Source.EnsureBareRepository(
				matches[2],
				path.Join(repository.HomeDirectory, repository.BareDirectory),
				path.Join(repository.HomeDirectory, repository.WorkingDirectory),
			); err != nil {
				fmt.Println(err)
			}
		} else {
			writer.Write([]byte("Unable to handle request"))
		}
	}
	repository.server.ServeHTTP(writer, request)
}
