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

type ArchiveSource struct {
	Host     string
	Versions []string
	Archive  string
}

type GitSource struct {
	Archive   *ArchiveSource
	Protocols []string
}

func (source *GitSource) EnsureBareRepository(name string, path string, cache string) error {
	workingPath := filepath.Join(cache, name)
	err := source.EnsureRepository(name, cache)
	if err != nil {
		return err
	}
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
	_, err = git.PlainClone(barePath, true, options)
	return err
}

func (source *GitSource) EnsureRepository(name string, path string) error {
	repositoyPath := filepath.Join(path, name)
	fmt.Println("Ensuring repository:", name, "from host:", source.Archive.Host, "on path:", path)
	if utils.VerifyPath(repositoyPath) {
		fmt.Println("Updating repository:", name)
		source.Archive.updateRepository(name, path)
	} else {
		options := &git.CloneOptions{
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		var err error
		for _, protocol := range source.Protocols {
			fmt.Println("Trying to clone repo with:", protocol)
			options.URL = fmt.Sprintf(selectUrlTemplate(protocol), source.Archive.Host, name)
			_, err = git.PlainClone(repositoyPath, false, options)
			if err == nil {
				fmt.Println("Repo clone successful using:", protocol)
				break
			}
		}
		if err != nil {
			fmt.Println("Unable to clone repo using git:", err)
			source.Archive.ensureTarballRepository(name, path)
		}
	}
	return nil
}

func selectUrlTemplate(protocol string) string {
	switch protocol {
	case "ssh":
		return protocol + "://git@%v/%v.git"
	}
	return protocol + "://%v/%v.git"
}

func (source *ArchiveSource) updateRepository(name string, path string) error {
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
		source.ensureTarballRepository(name, path)
	}
	return nil
}

func (source *ArchiveSource) ensureTarballRepository(name string, path string) error {
	for _, version := range source.Versions {
		url := fmt.Sprintf("https://%v/%v/%v/%v.tar.gz", source.Host, name, source.Archive, version)
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
