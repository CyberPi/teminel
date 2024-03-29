package proxy

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	proxyTls "source.cyberpi.de/go/teminel/core/proxy/tls"
	"source.cyberpi.de/go/teminel/utils"
	"source.cyberpi.de/go/teminel/utils/auth"
)

type Proxy struct {
	Target       *url.URL
	Credentials  *auth.Basic
	ReverseProxy *httputil.ReverseProxy
	TLSConfig    *proxyTls.Config
	Authenticate bool
}

func (proxy *Proxy) ListenAndServe(address string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.buildHandleFunc())
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	if proxy.TLSConfig == nil {
		fmt.Println("Info: Starting proxy:", server.Addr, "with target", proxy.Target)
		return server.ListenAndServe()
	} else {
		fmt.Println("Info: Starting proxy in TLS mode:", server.Addr, "with target", proxy.Target)
		server.TLSConfig = proxy.TLSConfig.Standard
		server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		return server.ListenAndServeTLS(proxy.TLSConfig.Cert, proxy.TLSConfig.Key)
	}
}

func (proxy *Proxy) forward(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	request.Host = proxy.Target.Host
	proxy.ReverseProxy.ServeHTTP(writer, request)
	end := time.Now()

	fmt.Println("Info:", request.Host,
		"Method:", request.Method,
		"Path:", request.URL.Path,
		"Time:", end.Format(time.RFC3339),
		"Latency:", end.Sub(start),
	)
}

func (proxy *Proxy) authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Add("Authorization", proxy.Credentials.FormatHeader())
		handler.ServeHTTP(writer, request)
	}
}

func (proxy *Proxy) authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		parts := strings.SplitN(request.Header.Get("Authorization"), " ", 2)
		if len(parts) != 2 {
			http.Error(writer, "Not authorized", 401)
			return
		}

		encoder, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(writer, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(encoder), ":", 2)
		if len(pair) != 2 {
			http.Error(writer, "Not authorized", 401)
			return
		}

		if pair[0] != proxy.Credentials.Name || pair[1] != proxy.Credentials.Password {
			http.Error(writer, "Not authorized", 401)
			return
		}

		handler.ServeHTTP(writer, request)
	}
}

func (proxy *Proxy) buildHandleFunc() http.HandlerFunc {
	if proxy.Credentials != nil {
		if proxy.Authenticate {
			fmt.Println("Proxy is authenticating at backend")
			return utils.Use(proxy.forward, proxy.authenticate)
		} else {
			fmt.Println("Setting up basic authentication")
			return utils.Use(proxy.forward, proxy.authorize)
		}
	} else {
		return proxy.forward
	}
}

func newReverseProxy(target *url.URL, insecure bool) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	if insecure {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return proxy
}
