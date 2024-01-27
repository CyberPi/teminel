package proxy

import (
	"crypto/tls"
	"fmt"
)

// for curl testing see https://unix.stackexchange.com/questions/208437/how-to-convert-ssl-ciphers-to-curl-format

var (
	tlsVersions = map[string]uint16{
		"TLS10": tls.VersionTLS10,
		"TLS11": tls.VersionTLS11,
		"TLS12": tls.VersionTLS12,
	}

	tlsCurves = map[string]tls.CurveID{
		"P256":   tls.CurveP256,
		"P384":   tls.CurveP384,
		"P521":   tls.CurveP521,
		"X25519": tls.X25519,
	}

	tlsCiphers = map[string]uint16{
		"RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
		"RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		"RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		"RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		"RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		"RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		"RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		"ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		"ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		"ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		"ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		"ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		"ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		"ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		"ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		"ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		"ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		"ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		"ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		"ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		"ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}
)

type TLSConfig struct {
	Min     string   `yaml:"min" json:"min"`
	Max     string   `yaml:"max" json:"max"`
	Curves  []string `yaml:"curves" json:"curves"`
	Ciphers []string `yaml:"ciphers" json:"ciphers"`
	Key     string   `yaml:"key" json:"key"`
	Cert    string   `yaml:"cert" json:"cert"`
}

func NewTLSConfig() *TLSConfig {
	return &TLSConfig{
		Min: "TLS10",
		Max: "TLS12",
		Curves: []string{
			"P521",
			"P384",
			"P256",
		},
		Ciphers: []string{
			"ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			"ECDHE_RSA_WITH_AES_256_CBC_SHA",
			"RSA_WITH_AES_256_GCM_SHA384",
			"RSA_WITH_AES_256_CBC_SHA",
		},
	}
}

func (config *TLSConfig) report() {
	fmt.Println("Info: Setting MIN TLS version",
		"TLSVersion:", config.Min,
	)
	fmt.Println("Info: Setting MAX TLS version",
		"TLSVersionID:", config.Max,
	)

	fmt.Println("Info: Setting Curve Preferences",
		"Curves:", config.Curves,
	)

	fmt.Println("Info: Setting Ciphers",
		"Ciphers:", config.Ciphers,
	)
}

func (config *TLSConfig) ToServerConfig() *tls.Config {
	curves := make([]tls.CurveID, len(config.Curves))
	for _, curveName := range config.Curves {
		curves = append(curves, tlsCurves[curveName])
	}
	ciphers := make([]uint16, len(config.Ciphers))
	for _, cipherName := range config.Ciphers {
		ciphers = append(ciphers, tlsCiphers[cipherName])
	}
	return &tls.Config{
		MinVersion:               tlsVersions[config.Min],
		MaxVersion:               tlsVersions[config.Max],
		PreferServerCipherSuites: true,
		CurvePreferences:         curves,
		CipherSuites:             ciphers,
	}
}
