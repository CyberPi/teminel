package tls

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
)

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

type readableConfig struct {
	Min     string   `yaml:"min" json:"min"`
	Max     string   `yaml:"max" json:"max"`
	Curves  []string `yaml:"curves" json:"curves"`
	Ciphers []string `yaml:"ciphers" json:"ciphers"`
	Key     string   `yaml:"key" json:"key"`
	Cert    string   `yaml:"cert" json:"cert"`
}

type Config struct {
	Key      string
	Cert     string
	Standard *tls.Config
}

func NewConfig() *Config {
	return &Config{
		Standard: &tls.Config{
			MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.CurveP384,
				tls.CurveP521,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}
}

func Unmarshal(data []byte, value *Config, unmarshal func([]byte, any) error) error {
	var readable readableConfig
	err := unmarshal(data, &readable)
	if err != nil {
		return err
	}
	curves := make([]tls.CurveID, len(readable.Curves))
	for _, curveName := range readable.Curves {
		curves = append(curves, tlsCurves[curveName])
	}
	ciphers := make([]uint16, len(readable.Ciphers))
	for _, cipherName := range readable.Ciphers {
		ciphers = append(ciphers, tlsCiphers[cipherName])
	}
	value.Standard = &tls.Config{
		MinVersion:               tlsVersions[readable.Min],
		MaxVersion:               tlsVersions[readable.Max],
		PreferServerCipherSuites: true,
		CurvePreferences:         curves,
		CipherSuites:             ciphers,
	}
	return err
}

func LoadJsonConfig(path string) (*Config, error) {
	if len(path) == 0 {
		return nil, nil
	}
	config := NewConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = Unmarshal(data, config, json.Unmarshal)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (config *Config) Report() {
	fmt.Println("Info: Setting MIN TLS version",
		"TLSVersion:", config.Standard.MinVersion,
	)
	fmt.Println("Info: Setting MAX TLS version",
		"TLSVersionID:", config.Standard.MaxVersion,
	)
	fmt.Println("Info: Setting Curve Preferences",
		"Curves:", config.Standard.CurvePreferences,
	)
	fmt.Println("Info: Setting Ciphers",
		"Ciphers:", config.Standard.CipherSuites,
	)
}
