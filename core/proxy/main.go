package proxy

import (
	"flag"
	"fmt"
	"net/url"

	"source.cyberpi.de/go/teminel/core/proxy/tls"
	"source.cyberpi.de/go/teminel/utils"
	"source.cyberpi.de/go/teminel/utils/auth"
)

func Main() {
	main()
}

func main() {
	userName := utils.EnsureEnv("USERNAME", "")
	flag.StringVar(&userName, "username", userName, "BasicAuth username to secure Proxy.")
	userPassword := utils.EnsureEnv("PASSWORD", "")
	flag.StringVar(&userPassword, "password", userPassword, "BasicAuth password to secure Proxy.")

	host := utils.EnsureEnv("TEMINEL_HOST", "0.0.0.0:8080")
	flag.StringVar(&host, "host", host, "Server IP address to bind to.")

	backend := utils.EnsureEnv("TEMINEL_BACKEND", "https://github.com:443")
	flag.StringVar(&backend, "backend", backend, "Backend server.")

	tlsPath := utils.EnsureEnv("TEMINEL_TLS", "")
	flag.StringVar(&tlsPath, "tls", tlsPath, "tls config file path.")

	insecure := false
	flag.BoolVar(&insecure, "insecure", insecure, "Skip backend tls verify.")

	authenticate := false
	flag.BoolVar(&authenticate, "auth", authenticate, "Authenticate at the backend only")

	flag.Parse()

	var user *auth.Basic
	if userName != "" && userPassword != "" {
		user = &auth.Basic{
			Name:     userName,
			Password: userPassword,
		}
	}

	targetUrl, err := url.Parse(backend)
	if err != nil {
		panic(fmt.Sprintln("Error: Unable to parse URL:", err))
	}

	var tlsConfig *tls.Config
	if len(tlsPath) == 0 {
		tlsConfig, err = tls.LoadJsonConfig(tlsPath)
		if err != nil {
			panic(err)
		}
	}

	proxy := &Proxy{
		Target:       targetUrl,
		Credentials:  user,
		ReverseProxy: newReverseProxy(targetUrl, insecure),
		TLSConfig:    tlsConfig,
		Authenticate: authenticate,
	}

	err = proxy.ListenAndServe(host)
	if err != nil {
		panic(err)
	}
}
