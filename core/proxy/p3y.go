package proxy

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"source.cyberpi.de/go/teminel/utils"
	"source.cyberpi.de/go/teminel/utils/auth"
)

type Proxy struct {
	Target       *url.URL
	Credentials  *auth.Basic
	reverseProxy *httputil.ReverseProxy
}

// forward requests
func (proxy *Proxy) forward(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	proxy.reverseProxy.ServeHTTP(writer, request)
	end := time.Now()

	fmt.Println("Info:", request.Host,
		"Method:", request.Method,
		"Path:", request.URL.Path,
		"Time:", end.Format(time.RFC3339),
		"Latency:", end.Sub(start),
	)
}

func (proxy *Proxy) authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if proxy.Credentials == nil {
			handler.ServeHTTP(writer, request)
			return
		}

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

// use
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// main entry point
func main() {
	backendEnv := utils.EnsureEnv("BACKEND", "http://github.com:80")
	usernameEnv := utils.EnsureEnv("USERNAME", "")
	passwordEnv := utils.EnsureEnv("PASSWORD", "")
	portEnv := utils.EnsureEnv("PORT", "80")
	ipEnv := utils.EnsureEnv("IP", "0.0.0.0")
	insecureEnv, _ := strconv.ParseBool(utils.EnsureEnv("INSECURE", "false"))
	tlsPathEnv := utils.EnsureEnv("TEMINEL_TLS", "")

	ip := flag.String("ip", ipEnv, "Server IP address to bind to.")
	port := flag.String("port", portEnv, "Server port.")
	backend := flag.String("backend", backendEnv, "backend server.")
	username := flag.String("username", usernameEnv, "BasicAuth username to secure Proxy.")
	password := flag.String("password", passwordEnv, "BasicAuth password to secure Proxy.")

	tlsPath := flag.String("tls", tlsPathEnv, "tls config file path.")
	insecure := flag.Bool("insecure", insecureEnv, "Skip backend tls verify.")

	flag.Parse()

	fmt.Println("Info", "Starting reverse proxy",
		"port", *port,
		"ip", *ip,
		"backend", *backend,
	)

	targetUrl, err := url.Parse(*backend)
	if err != nil {
		panic(fmt.Sprintln("Error: Unable to parse URL:", err))
	}

	// Proxy
	proxy := &Proxy{
		Target:       targetUrl,
		reverseProxy: httputil.NewSingleHostReverseProxy(targetUrl),
	}
	if *insecure {
		proxy.reverseProxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	mux := http.NewServeMux()

	if *username != "" {
		proxy.Credentials = &auth.Basic{
			Name:     *username,
			Password: *password,
		}
	}

	// server
	mux.HandleFunc("/", use(proxy.forward, proxy.authorize))

	srv := &http.Server{
		Addr:    *ip + ":" + *port,
		Handler: mux,
	}

	// If TLS is not specified serve the content unencrypted.
	if len(*tlsPath) == 0 {
		err = srv.ListenAndServe()
		if err != nil {
			panic(fmt.Sprintln("Error starting Proxy:", err))
		}
	} else {
		tlsConfig := NewTLSConfig()
		configData, err := os.ReadFile(*tlsPath)
		if err != nil {
			panic(fmt.Sprintln("Error on loading tlsConfig:", err))
		}
		err = json.Unmarshal(configData, tlsConfig)
		if err != nil {
			panic(fmt.Sprintln("Error on parsing tlsConfig:", err))
		}

		fmt.Println("Info: Starting Proxy in TLS mode.")
		srv.TLSConfig = tlsConfig.ToServerConfig()
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

		err = srv.ListenAndServeTLS(tlsConfig.Cert, tlsConfig.Key)
		if err != nil {
			panic(fmt.Sprintln("Error: starting proxyin TLS mode:", err))
		}
	}
}
