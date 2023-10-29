package tmux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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

func (tmux *Config) Load(path string) error {
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
