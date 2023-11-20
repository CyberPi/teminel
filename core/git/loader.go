package git

import (
	"fmt"
	"net/http"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/load"
)

type Loader struct {
	Target           string
	BareDirectory    string
	WorkingDirectory string
	Protocols        []string
	Versions         []string
	Archive          string
	server           *gitkit.Server
}

func (repository *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if err := load.CloneBare(
			repository.Protocols,
			repository.Target,
			matches[2],
			repository.BareDirectory,
			repository.WorkingDirectory,
			repository.Versions,
			repository.Archive,
		); err != nil {
			fmt.Println(err)
		}
	}
	repository.server.ServeHTTP(writer, request)
}
