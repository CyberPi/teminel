package neovim

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/sosedoff/gitkit"
)

const loaderDir = "/tmp/teminel/loader"

type Loader struct {
	server *gitkit.Server
}

func (cache *Loader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		fmt.Println("Host:", request.Host, "URL:", request.URL)
		matcher := regexp.MustCompile(`(^.+?)\.git.*$`)
		matches := matcher.FindStringSubmatch(request.URL.String())
		fmt.Println("Matches:", matches)
		http.Error(writer, "Could not read item", http.StatusNotFound)
		options := &git.CloneOptions{
			URL:          fmt.Sprintf("https://%v%v.git", "github.com", matches[1]),
			SingleBranch: true,
			Depth:        1,
			Tags:         git.NoTags,
		}
		repository, err := git.PlainClone(filepath.Join(loaderDir, matches[1]), false, options)
		if err != nil {
			fmt.Println(err)
		}
		repository.DeleteRemote("origin")
		repository.CreateRemote("origin", "http://localhost:9980")
	}
	cache.server.ServeHTTP(writer, request)
}
