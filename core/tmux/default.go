package tmux

import "source.cyberpi.de/go/teminel/core/git"

var defaultPath, _ = SelectConfig()
var Default Config = Config{
	Source: git.Default.Source,
	path:   defaultPath,
}
