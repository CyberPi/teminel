package git

import (
	"fmt"
	"net/http"
	"regexp"

	"source.cyberpi.de/go/teminel/load"
)

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

type Loader struct {
	Source           *load.GitSource
	BareDirectory    string
	WorkingDirectory string
}

func (repository *Loader) preload(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("New request:", request.URL.Path, "Method:", request.Method)
		if request.Method == http.MethodGet {
			matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
			if len(matches) < 2 {
				http.Error(writer, "Request is not an repository request", 400)
				return
			}
			if err := repository.Source.EnsureBareRepository(
				matches[2],
				repository.BareDirectory,
				repository.WorkingDirectory,
			); err != nil {
				http.Error(writer, fmt.Sprintf("Error on caching the request: %v", err), 500)
				return
			}
		}
		handler.ServeHTTP(writer, request)
	}
}
