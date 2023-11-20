package git

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/utils"
)

var repositoryMatcher = regexp.MustCompile(`^\/((.+?)\.git).*$`)

func Run(loader *Loader, port int) error {
	err := utils.EnsureDirectories(loader.BareDirectory)
	if err != nil {
		return err
	}
	mirrorConfig := gitkit.Config{
		Dir:        loader.BareDirectory,
		AutoCreate: true,
	}
	server := gitkit.New(mirrorConfig)
	if err := server.Setup(); err != nil {
		return err
	}
	loader.server = server
	http.Handle("/", loader)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
