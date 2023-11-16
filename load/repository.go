package load

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"source.cyberpi.de/go/teminel/utils"
)

const bareDirectory = "/tmp/teminel/bare"
const workingDirectory = "/tmp/teminel/working"

func CloneBare(protocol string, host string, path string, repositoryName string) error {
	bareRepository := repositoryName + ".git"
	barePath := filepath.Join(path, bareRepository)
	workingPath := filepath.Join(workingDirectory, repositoryName)
	// Check if working copy needs an update
	if utils.VerifyPath(workingPath) {
		updateTarballRepo(host, workingDirectory, repositoryName)
	}

	if !utils.VerifyPath(barePath) {
		fmt.Println("Loading repository:", repositoryName, "Using git clone from:", host)
		options := &git.CloneOptions{
			URL:          host + bareRepository,
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		_, err := git.PlainClone(barePath, true, options)
		if err != nil {
			fmt.Println("Try to get archive download due to error:", err)
			repository, err := git.PlainInit(workingPath, false)
			if err != nil {
				return err
			}
			worktree, err := repository.Worktree()
			if err != nil {
				return err
			}
			err = worktree.AddGlob(".")
			if err != nil {
				return err
			}

			options.URL = workingPath
			_, err = git.PlainClone(barePath, true, options)
			return err
		}
	}
	return nil
}

func updateTarballRepo(host string, path string, repositoryName string) error {
	url := fmt.Sprintf("https://%v/%v/archive/refs/heads/main.tar.gz", host, repositoryName)
	workingPath := filepath.Join(path, repositoryName)
	err := LoadTarball(url, workingPath)
	if err != nil {
		return err
	}
	repository, err := openOrInit(workingPath)
	if err != nil {
		return err
	}
	err = commit(repository)
	if err != nil {
		return err
	}
	return nil
}

func openOrInit(path string) (*git.Repository, error) {
	if utils.VerifyPath(path + "/.git") {
		return git.PlainOpen(path)
	} else {
		return git.PlainInit(path, false)
	}
}

func commit(repository *git.Repository) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}
	status, err := worktree.Status()
	if err != nil {
		return err
	}
	if !status.IsClean() {
		options := &git.AddOptions{All: true}
		worktree.AddWithOptions(options)
		_, err := worktree.Commit("New commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Auto Maton",
				Email: "teminel@cyberpi.de",
				When:  time.Now(),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
