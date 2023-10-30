package tmux

import "strings"

type Plugin struct {
	path string
	repo string
}

func (tmux *Plugin) Name() string {
	return strings.Split(tmux.path, "/")[1]
}
