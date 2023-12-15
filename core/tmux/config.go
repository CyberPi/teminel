package tmux

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"source.cyberpi.de/go/teminel/exec"
	"source.cyberpi.de/go/teminel/load"
	"source.cyberpi.de/go/teminel/utils"
)

type Config struct {
	path    string
	plugins []*Plugin
}

type Plugin struct {
	source string
	name   string
}

func (tmux *Config) Report() {
	fmt.Println(tmux.plugins)
	for _, plugin := range tmux.plugins {
		fmt.Println("repo:", plugin)
	}
}

var tpmPluginMatcher = regexp.MustCompile(`^set -g @(plugin(_source)) ["']?([\w/-]+?)["']?$`)

func (tmux *Config) Read(path string) error {
	tmux.path = filepath.Dir(path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	target := "github.com"
	for scanner.Scan() {
		line := scanner.Text()
		match := tpmPluginMatcher.FindStringSubmatch(line)
		if match != nil {
			switch match[1] {
			case "plugin_domain":
				target = match[3]
			default:
				tmux.plugins = append(tmux.plugins, &Plugin{
					source: target,
					name:   match[3],
				})
			}
		}
	}
	return nil
}

func (tmux *Config) Install() error {
	path := filepath.Join(tmux.path, "plugins")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.FileMode(0755))
	} else if err != nil {
		return err
	}
	for _, plugin := range tmux.plugins {

		load.EnsureRepository(
			[]string{"ssh", "https", "http"},
			plugin.source,
			plugin.name,
			path,
			[]string{"main", "master", "develop"},
			"archive/refs/heads",
		)
	}
	return nil
}

func (tmux *Config) Load() error {
	for _, plugin := range tmux.plugins {
		glob := fmt.Sprintf("%v/plugins/%v/*.tmux", tmux.path, path.Base(plugin.name))
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
