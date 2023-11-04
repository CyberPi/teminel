package neovim

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/sosedoff/gitkit"
	"source.cyberpi.de/go/teminel/utils"
)

const loaderDir = "/tmp/teminel/loader"

type Loader struct {
	server *gitkit.Server
}

var repositoryMatcher = regexp.MustCompile(`((^.+?)\.git).*$`)

func (cache *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		matches := repositoryMatcher.FindStringSubmatch(request.URL.String())
		repositoryPath := filepath.Join(mirrorDir, matches[1])
		if !utils.VerifyPath(repositoryPath) {
			fmt.Println("Loading repository:", matches[1], "Using git clone from:", request.Host)
			options := &git.CloneOptions{
				URL:          fmt.Sprintf("https://%v%v", "github.com", matches[1]),
				SingleBranch: true,
				Depth:        1,
				Tags:         git.NoTags,
			}
			_, err := git.PlainClone(repositoryPath, true, options)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	cache.server.ServeHTTP(writer, request)
}
