package neovim

import (
	"net/http"
	"os"

	"github.com/sosedoff/gitkit"
)

const mirrorDir = "/tmp/teminel/mirror/repos"

func Run() error {
	_, err := os.Stat(mirrorDir)
	if os.IsNotExist(err) {
		os.MkdirAll(mirrorDir, os.FileMode(0755))
	} else if err != nil {
		return err
	}
	_, err = os.Stat(loaderDir)
	if os.IsNotExist(err) {
		os.MkdirAll(loaderDir, os.FileMode(0755))
	} else if err != nil {
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
