package git

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sosedoff/gitkit"
	extFlag "source.cyberpi.de/go/teminel/flag"
	"source.cyberpi.de/go/teminel/load"
	"source.cyberpi.de/go/teminel/utils"
)

func Main() {
	main()
}

func main() {
	host := utils.EnsureEnv("TEMINEL_HOST", "0.0.0.0:8080")
	flag.StringVar(&host, "host", host, "Server IP address to bind to.")

	backend := utils.EnsureEnv("TEMINEL_BACKEND", "github.com")
	flag.StringVar(&backend, "backend", backend, "Backend server.")
	var versions extFlag.MultiFlag
	flag.Var(&versions, "version", "Array of branches to check for.")
	versions.Default("main", "master", "develop")
	archive := utils.EnsureEnv("TEMINEL_ARCHIVE", "archive/refs/heads")
	flag.StringVar(&archive, "archive", archive, "Path to tar archives")

	var protocols extFlag.MultiFlag
	flag.Var(&protocols, "protocol", "Protocols to use to clone git repo")
	protocols.Default("ssh", "https", "http")

	home := utils.EnsureEnv("TEMINEL_HOME", "/var/lib/teminel")
	flag.StringVar(&home, "home", home, "Main dir to store teminel caches and config.")
	bare := utils.EnsureEnv("TEMINEL_BARE", "bare")
	flag.StringVar(&bare, "bare", bare, "Sub dir to store teminel bare repo data.")
	working := utils.EnsureEnv("TEMINEL_WORKING", "working")
	flag.StringVar(&working, "working", working, "Sub dir to store teminel repo data.")

	flag.Parse()

	loader := Loader{
		Source: &load.GitSource{
			Archive: &load.ArchiveSource{
				Host:     backend,
				Versions: versions,
				Archive:  archive,
			},
			Protocols: protocols,
		},
		BareDirectory:    filepath.Join(home, bare),
		WorkingDirectory: filepath.Join(home, working),
	}

	err := utils.EnsureDirectories(loader.BareDirectory)
	if err != nil {
		panic(err)
	}
	mirrorConfig := gitkit.Config{
		Dir:        loader.BareDirectory,
		AutoCreate: true,
	}
	server := gitkit.New(mirrorConfig)
	if err := server.Setup(); err != nil {
		panic(err)
	}
	loader.server = server
	http.Handle("/", &loader)
	fmt.Println("Git mirror proxy started at:", host)
	panic(http.ListenAndServe(host, nil))
}
