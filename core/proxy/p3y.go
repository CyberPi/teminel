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

var Version = "0.0.1"

// Proxy defines the Proxy handler see NewProx()
type Proxy struct {
	Target       *url.URL
	Credentials  *auth.Basic
	reverseProxy *httputil.ReverseProxy
}

// handle requests
func (proxy *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	proxy.reverseProxy.ServeHTTP(w, r)
	end := time.Now()

	fmt.Println("Info:", r.Host,
		"Method:", r.Method,
		"Path:", r.URL.Path,
		"Time:", end.Format(time.RFC3339),
		"Latency:", end.Sub(start),
	)
}

// basicAuth
func (proxy *Proxy) basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if proxy.Credentials == nil {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != proxy.Credentials.Name || pair[1] != proxy.Credentials.Password {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
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

	version := flag.Bool("version", false, "Display version.")
	flag.Parse()

	if *version {
		fmt.Println("Version:", Version)
		os.Exit(1)
	}

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
	mux.HandleFunc("/", use(proxy.handle, proxy.basicAuth))

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
