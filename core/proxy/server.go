package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"source.cyberpi.de/go/teminel/utils/auth"
)

type Proxy struct {
	Target *url.URL
	Auth   *auth.Basic
	server *httputil.ReverseProxy
}

func Run(proxy *Proxy, port int) error {
	server := httputil.NewSingleHostReverseProxy(proxy.Target)
	proxy.server = server

	if proxy.Auth != nil {
		http.HandleFunc("/", proxy.ForwardWithBasicAuth())
	} else {
		http.HandleFunc("/", proxy.Forward())
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func (proxy *Proxy) ForwardWithBasicAuth() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Add("Authorization", proxy.Auth.FormatHeader())
		proxy.server.ServeHTTP(writer, request)
	}
}

func (proxy *Proxy) Forward() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		proxy.server.ServeHTTP(writer, request)
	}
}