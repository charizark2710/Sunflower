package handler

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

var tlsConfig *tls.Config

func ReadRootCACert() error {
	caCert, err := os.ReadFile(os.Getenv("ROOT_CA_CERT_PATH"))
	if err != nil {
		return err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig = &tls.Config{
		RootCAs: caPool,
	}
	return nil
}

func GetTlsConfig() *tls.Config {
	return tlsConfig
}
