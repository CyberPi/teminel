package tmux

import (
	"fmt"
	"os"
	"path/filepath"

	"source.cyberpi.de/go/teminel/utils"
)

func Run() error {
	_, isTmuxSession := os.LookupEnv("TMUX")
	configFile, err := SelectConfig()
	if err != nil {
		return err
	}
	config := &Config{}
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
