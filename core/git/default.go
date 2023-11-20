package git

var Default Loader = Loader{
	Target:           "github.com",
	BareDirectory:    "/tmp/teminel/bare",
	WorkingDirectory: "/tmp/teminel/working",
	Protocols:        []string{"ssh, https, http"},
	Versions:         []string{"main", "master", "develop"},
	Archive:          "archive/refs/heads",
}
