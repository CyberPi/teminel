package proxy

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// for curl testing see https://unix.stackexchange.com/questions/208437/how-to-convert-ssl-ciphers-to-curl-format

var (
	TLSVersions = map[string]uint16{
		"VersionTLS10": tls.VersionTLS10,
		"VersionTLS11": tls.VersionTLS11,
		"VersionTLS12": tls.VersionTLS12,
	}

	Curves = map[string]tls.CurveID{
		"CurveP256": tls.CurveP256,
		"CurveP384": tls.CurveP384,
		"CurveP521": tls.CurveP521,
		"X25519":    tls.X25519,
	}

	Ciphers = map[string]uint16{
		"TLS_RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		"TLS_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}
)

type TLSPreferences struct {
	Min              string   `yaml:"min"`
	Max              string   `yaml:"max"`
	CurvePreferences []string `yaml:"curvePreferences"`
	Ciphers          []string `yaml:"ciphers"`
}

func GenericTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS10,
		MaxVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}

// NewTLSCfgFromYaml
func NewTLSCfgFromJson(filename string) (*tls.Config, error) {

	ymlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tlsPreferences := TLSPreferences{}

	err = json.Unmarshal([]byte(ymlData), &tlsPreferences)
	if err != nil {
		return nil, err
	}

	fmt.Println("Info: Setting MIN TLS version",
		"TLSVersionName:", tlsPreferences.Min,
		"TLSVersionID:", TLSVersions[tlsPreferences.Min],
	)

	fmt.Println("Info: Setting MAX TLS version",
		"TLSVersionName:", tlsPreferences.Max,
		"TLSVersionID:", TLSVersions[tlsPreferences.Max],
	)

	fmt.Println("Info: Setting Curve Preferences",
		"Curves:", tlsPreferences.CurvePreferences,
	)

	fmt.Println("Info: Setting Ciphers",
		"Ciphers:", tlsPreferences.Ciphers,
	)

	curveIDs := make([]tls.CurveID, 0)
	for _, curveName := range tlsPreferences.CurvePreferences {
		curveIDs = append(curveIDs, Curves[curveName])
	}

	cipher := make([]uint16, 0)
	for _, cipherName := range tlsPreferences.Ciphers {
		cipher = append(cipher, Ciphers[cipherName])
	}

	tlsCfg := &tls.Config{
		MinVersion:               TLSVersions[tlsPreferences.Min],
		MaxVersion:               TLSVersions[tlsPreferences.Max],
		CurvePreferences:         curveIDs,
		PreferServerCipherSuites: true,
		CipherSuites:             cipher,
	}

	return tlsCfg, nil
}
