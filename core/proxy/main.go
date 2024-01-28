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
		Credentials:  auth.NewBasic(*username, *password),
		ReverseProxy: newReverseProxy(targetUrl, *insecure),
		TLSConfig:    tlsConfig,
	}

	err = proxy.ListenAndServe(*ip + ":" + *port)
	if err != nil {
		panic(err)
	}
}
