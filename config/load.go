package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"source.cyberpi.de/go/teminel/utils"
)

type Config struct {
	RootDir string
	Plugins []string
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	matcher := regexp.MustCompile(`set -g @plugin ["']([A-z-/]+)["']`)
	result := Config{}
	for scanner.Scan() {
		line := scanner.Text()
		match := matcher.FindStringSubmatch(line)
		if match != nil {
			result.Plugins = append(result.Plugins, match[1])
		}
	}
	return &result, nil
}

func SelectConfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return homeDir, err
	}
	possibilities := []string{
		".tmux.conf",
		".tmux/tmux.conf",
		".config/tmux/tmux.conf",
	}
	for _, possibility := range possibilities {
		toCheck := fmt.Sprintf("%v/%v", homeDir, possibility)
		fmt.Printf("to check: %v\n", toCheck)
		if utils.CheckFileExists(toCheck) {
			return toCheck, nil
		}
	}
	return homeDir, fmt.Errorf("No possible configuration found")
}
