package tmux

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"source.cyberpi.de/go/teminel/utils"
)

type Config struct {
	path    string
	plugins []*Plugin
}

type Plugin struct {
	path string
	repo string
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
	matcher := regexp.MustCompile(`(PLUGIN ["']([\w\/-]+?)["'])|(REPO ["']([\w\/-]+?)["'])`)
	for scanner.Scan() {
		line := scanner.Text()
		match := matcher.FindStringSubmatch(line)
		if match != nil {
			if match[4] == "" {
				match[4] = "github.com"
			}
			toAppend := &Plugin{
				path: match[2],
				repo: match[4],
			}
			tmux.plugins = append(tmux.plugins, toAppend)
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
		pluginName := strings.Split(plugin.path, "/")[1]
		installPath := filepath.Join(pluginPath, pluginName)
		_, err := os.Stat(installPath)
		if os.IsNotExist(err) {
			fmt.Println("Installing plugin", plugin, "in", installPath)
			options := &git.CloneOptions{
				URL:          fmt.Sprintf("https://%v/%v.git", plugin.repo, plugin.path),
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
	}
	return nil
}

func (tmux *Config) Load() error {
	for _, plugin := range tmux.plugins {
		pluginName := strings.Split(plugin.path, "/")[1]
		glob := fmt.Sprintf("%v/plugins/%v/*.tmux", tmux.path, pluginName)
		toLoad, err := filepath.Glob(glob)
		if err != nil {
			return err
		}
		for _, item := range toLoad {
			command := exec.Command("sh", "-c", item)
			stdout, err := command.Output()
			if err != nil {
				return err
			}
			fmt.Printf("Load plugin: %v\n", item)
			for _, line := range stdout {
				fmt.Println(line)
			}
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
		toCheck := fmt.Sprintf("%v/%v", homeDir, possibility)
		if utils.CheckFileExists(toCheck) {
			return toCheck, nil
		}
	}
	return homeDir, fmt.Errorf("No possible configuration found")
}
