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
	BareDirectory:    "var/lib/teminel/bare",
	WorkingDirectory: "var/lib/teminel/working",
}
