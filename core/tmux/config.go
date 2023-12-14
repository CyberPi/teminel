package tmux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"source.cyberpi.de/go/teminel/exec"
)

type Config struct {
	path    string
	plugins []*Plugin
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
	for scanner.Scan() {
		toAppend := ParsePlugin(scanner.Text())
		if toAppend != nil {
			tmux.plugins = append(tmux.plugins, toAppend)
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
		plugin.Install(path)
	}
	return nil
}

func (tmux *Config) Load() error {
	for _, plugin := range tmux.plugins {
		glob := fmt.Sprintf("%v/plugins/%v/*.tmux", tmux.path, plugin.Name())
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
