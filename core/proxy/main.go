package proxy

import (
	"flag"
	"fmt"
	"net/url"
	"strconv"

	"source.cyberpi.de/go/teminel/utils"
	"source.cyberpi.de/go/teminel/utils/auth"
)

func main() {
	user := &auth.Basic{
		Name:     utils.EnsureEnv("USERNAME", ""),
		Password: utils.EnsureEnv("PASSWORD", ""),
	}
	user.Name = *flag.String("username", user.Name, "BasicAuth username to secure Proxy.")
	user.Password = *flag.String("password", user.Password, "BasicAuth password to secure Proxy.")
	if user.Name == "" || user.Password == "" {
		user = nil
	}

	hostEnv := utils.EnsureEnv("host", "0.0.0.0:80")
	host := *flag.String("host", hostEnv, "Server IP address to bind to.")

	backendEnv := utils.EnsureEnv("BACKEND", "http://github.com:80")
	backend := flag.String("backend", backendEnv, "backend server.")

	tlsPathEnv := utils.EnsureEnv("TEMINEL_TLS", "")
	tlsPath := flag.String("tls", tlsPathEnv, "tls config file path.")

	insecure := *flag.Bool("insecure", false, "Skip backend tls verify.")

	flag.Parse()

	targetUrl, err := url.Parse(*backend)
	if err != nil {
		panic(fmt.Sprintln("Error: Unable to parse URL:", err))
	}

	tlsConfig, err := loadTLSConfig(*tlsPath)
	if err != nil {
		panic(err)
	}

	proxy := &Proxy{
		Target:       targetUrl,
		Credentials:  user,
		ReverseProxy: newReverseProxy(targetUrl, insecure),
		TLSConfig:    tlsConfig,
	}

	err = proxy.ListenAndServe(host)
	if err != nil {
		panic(err)
	}
}
