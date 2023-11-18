package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"source.cyberpi.de/go/teminel/utils/auth"
)

type Config struct {
	Target *url.URL
	Port   int
	Auth   *auth.Basic
}

func Run(config *Config) error {
	server := httputil.NewSingleHostReverseProxy(config.Target)

	if config.Auth != nil {
		http.HandleFunc("/", WithBasicAuth(server, config.Auth))
	} else {
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			server.ServeHTTP(writer, request)
		})
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
}

func WithBasicAuth(server *httputil.ReverseProxy, basicAuth *auth.Basic) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Add("Authorization", basicAuth.FormatHeader())
		server.ServeHTTP(writer, request)
	}
}
