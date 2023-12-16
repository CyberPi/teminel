package tmux

import (
	"os"
)

func Run() error {
	_, isTmuxSession := os.LookupEnv("TMUX")
	err := Default.Install()
	if err != nil {
		return err
	}
	if isTmuxSession {
		return Default.Load()
	}
	return nil
}
