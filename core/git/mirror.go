package neovim

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/load"
	"source.cyberpi.de/go/teminel/utils"
)

const bareDirectory = "/tmp/teminel/bare"
const workingDirectory = "/tmp/teminel/working"

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

func Run() error {
	err := utils.EnsureDirectories(bareDirectory)
	if err != nil {
		return err
	}
	config := gitkit.Config{
		Dir:        bareDirectory,
		AutoCreate: true,
	}
	mirror := gitkit.New(config)
	if err := mirror.Setup(); err != nil {
		return err
	}
	loader := &Loader{
		mirror: mirror,
	}
	http.Handle("/", loader)

	return http.ListenAndServe(fmt.Sprintf(":%v", 9980), nil)
}

type Loader struct {
	mirror *gitkit.Server
}

func (cache *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		if err := load.CloneBare("github.com", matches[2], bareDirectory, workingDirectory); err != nil {
			fmt.Println(err)
		}
	}
	cache.mirror.ServeHTTP(writer, request)
}
