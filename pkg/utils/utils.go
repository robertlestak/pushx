package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func TlsConfig(enableTLS *bool, TLSInsecure *bool, TLSCA *string, TLSCert *string, TLSKey *string) (*tls.Config, error) {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "tlsConfig",
	})
	l.Debug("Creating TLS config")
	tc := &tls.Config{}
	if enableTLS != nil && *enableTLS {
		l.Debug("Enabling TLS")
		if TLSInsecure != nil && *TLSInsecure {
			l.Debug("Enabling TLS insecure")
			tc.InsecureSkipVerify = true
		}
		if TLSCA != nil && *TLSCA != "" {
			l.Debug("Enabling TLS CA")
			caCert, err := ioutil.ReadFile(*TLSCA)
			if err != nil {
				l.Errorf("%+v", err)
				return tc, err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tc.RootCAs = caCertPool
		}
		if TLSCert != nil && *TLSCert != "" {
			l.Debug("Enabling TLS cert")
			cert, err := tls.LoadX509KeyPair(*TLSCert, *TLSKey)
			if err != nil {
				l.Errorf("%+v", err)
				return tc, err
			}
			tc.Certificates = []tls.Certificate{cert}
		}
	}
	l.Debug("Created TLS config")
	return tc, nil
}
