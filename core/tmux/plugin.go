package tmux

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

const defaultSrc = "github.com"

type Plugin struct {
	path string
	repo string
}

func (tmux *Plugin) Name() string {
	return strings.Split(tmux.path, "/")[1]
}

func (tmux *Plugin) Install(path string) error {
	installPath := filepath.Join(path, tmux.Name())
	_, err := os.Stat(installPath)
	if os.IsNotExist(err) {
		fmt.Println("Installing plugin", tmux.path, "in", installPath)
		options := &git.CloneOptions{
			URL:          fmt.Sprintf("https://%v/%v.git", tmux.repo, tmux.path),
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		_, err = git.PlainClone(installPath, false, options)
	}
	return err
}

var pluginMatcher = regexp.MustCompile(`(PLUGIN ["']([\w\/-]+?)["'])|(REPO ["']([\w\/-]+?)["'])`)

// set -g @plugin 'tmux-plugins/tmux-sensible'
var tpmPluginRepoMatcher = regexp.MustCompile(`^set -g @plugin_repo ["']?([\w/-]+?)["']?$`)
var tpmPluginMatcher = regexp.MustCompile(`^set -g @plugin ["']?([\w/-]+?)["']?$`)

func ParsePlugin(line string) *Plugin {
	match := pluginMatcher.FindStringSubmatch(line)
	if match != nil {
		if match[4] == "" {
			match[4] = defaultSrc
		}
		return &Plugin{
			path: match[2],
			repo: match[4],
		}
	}
	return nil
}
