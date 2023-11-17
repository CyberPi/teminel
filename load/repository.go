package load

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"source.cyberpi.de/go/teminel/utils"
)

func CloneBare(host string, name string, path string, cache string) error {
	workingPath := filepath.Join(cache, name)
	ensureRepository(host, name, cache)

	bareRepository := name + ".git"
	barePath := filepath.Join(path, bareRepository)
	if utils.VerifyPath(barePath) {
		os.RemoveAll(barePath)
	}
	options := &git.CloneOptions{
		URL:          workingPath,
		SingleBranch: true,
		Depth:        1,
		Tags:         git.NoTags,
	}
	_, err := git.PlainClone(barePath, true, options)
	return err
}

func ensureRepository(host string, name string, path string) error {
	repositoyPath := filepath.Join(path, name)
	fmt.Println("Ensuring repository:", name, "from host:", host, "on path:", path)
	if utils.VerifyPath(repositoyPath) {
		updateRepository(host, name, path)
	} else {
		urls := []string{
			fmt.Sprintf("ssh://git@%v/%v.git", host, name),
			fmt.Sprintf("https://%v/%v.git", host, name),
		}
		options := &git.CloneOptions{
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		var err error
		for _, url := range urls {
			fmt.Println("Trying to clone repo from:", url)
			options.URL = url
			_, err = git.PlainClone(repositoyPath, false, options)
			if err == nil {
				fmt.Println("Repo cloned using git")
				break
			}
		}
		if err != nil {
			fmt.Println("Unable to clone repo using git:", err)
			ensureTarballRepository(host, name, path)
		}
	}
	return nil
}

func updateRepository(host string, name string, path string) error {
	workingPath := filepath.Join(path, name)
	repository, err := openOrInit(workingPath)
	remotes, err := repository.Remotes()
	if err != nil {
		return nil
	}
	if len(remotes) != 0 {
		worktree, err := repository.Worktree()
		if err != nil {
			return err
		}
		options := &git.PullOptions{
			RemoteName: "origin",
		}
		err = worktree.Pull(options)
		if err == git.NoErrAlreadyUpToDate {

		} else if err != nil {
			return nil
		}
	} else {
		ensureTarballRepository(host, name, path)
	}
	return nil
}

func ensureTarballRepository(host string, name string, path string) error {
	branches := []string{
		"main",
		"master",
		"develop",
	}
	for _, branch := range branches {
		url := fmt.Sprintf("https://%v/%v/archive/refs/heads/%v.tar.gz", host, name, branch)
		workingPath := filepath.Join(path, name)
		err := LoadTarball(url, workingPath)
		if err != nil {
			fmt.Println("Tarball failed to load:", err)
			continue
		}
		repository, err := openOrInit(workingPath)
		if err != nil {
			return err
		}
		err = commit(repository)
		if err != nil {
			return err
		}
		break
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
