package git

import "source.cyberpi.de/go/teminel/load"

var Default Loader = Loader{
	Source: &load.GitSource{
		Archive: &load.ArchiveSource{
			Host:     "github.com",
			Versions: []string{"main", "master", "develop"},
			Archive:  "archive/refs/heads",
		},
		Protocols: []string{"ssh", "https", "http"},
	},
	BareDirectory:    "/tmp/teminel/bare",
	WorkingDirectory: "/tmp/teminel/working",
}
