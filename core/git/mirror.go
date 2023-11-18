package neovim

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/load"
	"source.cyberpi.de/go/teminel/utils"
)

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

type Config struct {
	Port             int
	BareDirectory    string
	workingDirectory string
	GitProtocols     []string
	GitBranches      []string
}

func Run(config *Config) error {
	err := utils.EnsureDirectories(config.BareDirectory)
	if err != nil {
		return err
	}
	mirrorConfig := gitkit.Config{
		Dir:        config.BareDirectory,
		AutoCreate: true,
	}
	server := gitkit.New(mirrorConfig)
	if err := server.Setup(); err != nil {
		return err
	}
	loader := &loader{
		server: server,
	}
	http.Handle("/", loader)

	return http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
}

type loader struct {
	server           *gitkit.Server
	target           string
	bareDirectory    string
	workingDirectory string
}

func (repository *loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if err := load.CloneBare(
			repository.target,
			matches[2],
			repository.bareDirectory,
			repository.workingDirectory,
		); err != nil {
			fmt.Println(err)
		}
	}
	repository.server.ServeHTTP(writer, request)
}
