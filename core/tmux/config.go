package tmux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"source.cyberpi.de/go/teminel/exec"
	"source.cyberpi.de/go/teminel/load"
	"source.cyberpi.de/go/teminel/utils"
)

type Config struct {
	Source  *load.GitSource
	path    string
	plugins []*Plugin
}

type Plugin struct {
	host string
	name string
}

var tpmPluginMatcher = regexp.MustCompile(`^set -g @(plugin(_host|_path)?) ["']?([\w/-]+?)["']?$`)

func (tmux *Config) Read(configPath string) error {
	tmux.path = filepath.Join(filepath.Dir(configPath), "plugins")
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	host := tmux.Source.Archive.Host
	for scanner.Scan() {
		line := scanner.Text()
		match := tpmPluginMatcher.FindStringSubmatch(line)
		if match != nil {
			switch match[1] {
			case "plugin_host":
				host = match[3]
			case "plugin_path":
				tmux.path = match[3]
			default:
				tmux.plugins = append(tmux.plugins, &Plugin{
					host: host,
					name: match[3],
				})
			}
		}
	}
	return nil
}

func (tmux *Config) Install() error {
	_, err := os.Stat(tmux.path)
	if os.IsNotExist(err) {
		os.Mkdir(tmux.path, os.FileMode(0755))
	} else if err != nil {
		return err
	}
	for _, plugin := range tmux.plugins {
		tmux.Source.Archive.Host = plugin.host
		tmux.Source.EnsureRepository(plugin.name, tmux.path)
	}
	return nil
}

func (tmux *Config) Load() error {
	for _, plugin := range tmux.plugins {
		fmt.Println("Start plugin:", plugin.name)
		glob := fmt.Sprintf("%v/%v/*.tmux", tmux.path, plugin.name)
		toLoad, err := filepath.Glob(glob)
		if err != nil {
			return err
		}
		for _, item := range toLoad {
			exec.Shell(item)
		}
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
		toCheck := filepath.Join(homeDir, possibility)
		if utils.VerifyPath(toCheck) {
			return toCheck, nil
		}
	}
	return homeDir, fmt.Errorf("No possible configuration found")
}
