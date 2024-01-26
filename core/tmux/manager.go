package tmux

import (
	"os"

	"source.cyberpi.de/go/teminel/core/git"
)

func Run() error {
	_, isTmuxSession := os.LookupEnv("TMUX")
	configFile, err := SelectConfig()
	if err != nil {
		return err
	}
	config := Config{
		Source: git.Default.Source,
	}
	err = config.Read(configFile)
	if err != nil {
		return err
	}
	err = config.Install()
	if err != nil {
		return err
	}
	if isTmuxSession {
		return config.Load()
	}
	return nil
}
