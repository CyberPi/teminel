package tmux

import (
	"flag"
	"os"

	extFlag "source.cyberpi.de/go/teminel/flag"
	"source.cyberpi.de/go/teminel/load"
	"source.cyberpi.de/go/teminel/utils"
)

func Main() {
	main()
}

func main() {
	configFile, err := SelectConfig()
	if err != nil {
		panic(err)
	}
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

	install := false
	flag.BoolVar(&install, "install", install, "Runs the installation of plufins")

	flag.Parse()

	config := Config{
		Source: &load.GitSource{
			Archive: &load.ArchiveSource{
				Host:        backend,
				Versions:    versions,
				Archive:     archive,
				UseBaseName: true,
			},
			Protocols: protocols,
		},
	}
	err = config.Read(configFile)
	if err != nil {
		panic(err)
	}
	if install {
		err = config.Install()
		if err != nil {
			panic(err)
		}
	}
	_, isTmuxSession := os.LookupEnv("TMUX")
	if isTmuxSession {
		err := config.Load()
		if err != nil {
			panic(err)
		}
	}
}
