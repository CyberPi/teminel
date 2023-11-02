package neovim

import (
	"net/http"

	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/utils"
)

const mirrorDir = "/tmp/teminel/mirror/repos"

func Run() error {
	err := utils.EnsureDirectories(mirrorDir, loaderDir)
	if err != nil {
		return err
	}
	hooks := &gitkit.HookScripts{
		PreReceive: `echo "This is just a mirror. You should notpush here!"`,
	}
	config := gitkit.Config{
		Dir:        mirrorDir,
		AutoCreate: true,
		AutoHooks:  true,
		Hooks:      hooks,
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
