package proxy

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var Version = "0.0.1"

// BasicCredentials supplied by flags
type BasicCredentials struct {
	Username string
	Password string
}

// Proxy defines the Proxy handler see NewProx()
type Proxy struct {
	Target      *url.URL
	Proxy       *httputil.ReverseProxy
	Credentials *BasicCredentials
}

// handle requests
func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	p.Proxy.ServeHTTP(w, r)
	end := time.Now()

	fmt.Println("Info:", r.Host,
		"Method:", r.Method,
		"Path:", r.URL.Path,
		"Time:", end.Format(time.RFC3339),
		"Latency:", end.Sub(start),
	)
}

// basicAuth
func (p *Proxy) basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if p.Credentials == nil {
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

		if pair[0] != p.Credentials.Username || pair[1] != p.Credentials.Password {
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
	backendEnv := getEnv("BACKEND", "http://github.com:80")
	usernameEnv := getEnv("USERNAME", "")
	passwordEnv := getEnv("PASSWORD", "")
	portEnv := getEnv("PORT", "80")
	ipEnv := getEnv("IP", "0.0.0.0")
	crtEnv := getEnv("CRT", "./CA.crt")
	keyEnv := getEnv("KEY", "./CA.key")
	insecureEnv, _ := strconv.ParseBool(getEnv("INSECURE", "false"))
	tlsCfgFileEnv := getEnv("TLSCFG", "")
	tlsEnv, _ := strconv.ParseBool(getEnv("TLS", "false"))

	ip := flag.String("ip", ipEnv, "Server IP address to bind to.")
	port := flag.String("port", portEnv, "Server port.")
	backend := flag.String("backend", backendEnv, "backend server.")
	username := flag.String("username", usernameEnv, "BasicAuth username to secure Proxy.")
	password := flag.String("password", passwordEnv, "BasicAuth password to secure Proxy.")
	srvtls := flag.Bool("tls", tlsEnv, "TLS Support (requires crt and key)")
	tlsCfgFile := flag.String("tlsCfg", tlsCfgFileEnv, "tls config file path.")
	crt := flag.String("crt", crtEnv, "Path to cert. (enable --tls)")
	key := flag.String("key", keyEnv, "Path to private key. (enable --tls")
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

	pxy := httputil.NewSingleHostReverseProxy(targetUrl)

	if *insecure {
		pxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Proxy
	proxy := &Proxy{
		Target: targetUrl,
		Proxy:  pxy,
	}

	mux := http.NewServeMux()

	if *username != "" {
		proxy.Credentials = &BasicCredentials{
			Username: *username,
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
	if *srvtls != true {
		err = srv.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting Proxy: %s\n", err.Error())
		}
		os.Exit(0)
	}

	// Get a generic TLS configuration
	tlsCfg := GenericTLSConfig()
	if *tlsCfgFile == "" {
		fmt.Println("Warn: No TLS configuration specified, using default.")
	}

	if *tlsCfgFile != "" {
		fmt.Println("Info: Loading TLS configuration from " + *tlsCfgFile)
		tlsCfg, err = NewTLSCfgFromJson(*tlsCfgFile)
		if err != nil {
			fmt.Println("Error: configuring TLS:", err)
			os.Exit(0)
		}
	}

	fmt.Println("Info: Starting Proxy in TLS mode.")

	srv.TLSConfig = tlsCfg
	srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	err = srv.ListenAndServeTLS(*crt, *key)
	if err != nil {
		fmt.Println("Error: starting proxyin TLS mode:", err)
	}

}

// getEnv gets an environment variable or sets a default if
// one does not exist.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
