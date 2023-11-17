package neovim

import (
	"net/http"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/utils"
)

const bareDirectory = "/tmp/teminel/bare"

func Run() error {
	err := utils.EnsureDirectories(bareDirectory)
	if err != nil {
		return err
	}
	config := gitkit.Config{
		Dir:        bareDirectory,
		AutoCreate: true,
	}
	middleware := gitkit.New(config)
	if err := middleware.Setup(); err != nil {
		return err
	}
	loader := &Loader{
		server: middleware,
	}
	http.Handle("/", loader)

	return http.ListenAndServe(":9980", nil)
}
