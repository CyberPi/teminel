package git

import (
	"fmt"
	"net/http"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/load"
)

type Loader struct {
	Source           *load.GitSource
	BareDirectory    string
	WorkingDirectory string
	server           *gitkit.Server
}

func (repository *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if err := repository.Source.EnsureBareRepository(
			matches[2],
			repository.BareDirectory,
			repository.WorkingDirectory,
		); err != nil {
			fmt.Println(err)
		}
	}
	repository.server.ServeHTTP(writer, request)
}
