package tmux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"source.cyberpi.de/go/teminel/utils"
)

type Config struct {
	path    string
	plugins []string
}

func (tmux *Config) Report() {
	fmt.Println(tmux.plugins)
	for _, plugin := range tmux.plugins {
		fmt.Println("repo:", plugin)
	}
}

func (tmux *Config) Read(path string) error {
	tmux.path = filepath.Dir(path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	matcher := regexp.MustCompile(`^#PLUGIN ["'](.+?)["']$`)
	for scanner.Scan() {
		line := scanner.Text()
		match := matcher.FindStringSubmatch(line)
		if match != nil {
			tmux.plugins = append(tmux.plugins, match[1])
		}
	}
	return nil
}

func (tmux *Config) Install() error {
	pluginPath := filepath.Join(tmux.path, "plugins")
	_, err := os.Stat(pluginPath)
	if os.IsNotExist(err) {
		os.Mkdir(pluginPath, os.FileMode(0755))
	} else if err != nil {
		return err
	}
	for _, plugin := range tmux.plugins {
		pluginName := strings.Split(plugin, "/")[1]
		installPath := filepath.Join(pluginPath, pluginName)
		_, err := os.Stat(installPath)
		if os.IsNotExist(err) {
			fmt.Println("Installing plugin", plugin, "in", installPath)
			options := &git.CloneOptions{
				URL:          fmt.Sprintf("https://github.com/%v.git", plugin),
				RemoteName:   "origin",
				SingleBranch: true,
				Mirror:       false,
				NoCheckout:   false,
				Depth:        1,
				Tags:         git.NoTags,
			}
			_, err := git.PlainClone(installPath, false, options)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		fmt.Println("Plugin", pluginName, "installed in", installPath)
	}
	return nil
}

func SelectConfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return homeDir, err
	}
	possibilities := []string{
		".tmux.conf",
		".config/tmux/tmux.conf",
	}
	for _, possibility := range possibilities {
		toCheck := fmt.Sprintf("%v/%v", homeDir, possibility)
		if utils.CheckFileExists(toCheck) {
			return toCheck, nil
		}
	}
	return homeDir, fmt.Errorf("No possible configuration found")
}
