package crypt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
)

// TLSServerConfig provides config with cert and key
func TLSServerConfig() *tls.Config {
	caCertPEM, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(os.Args[4], os.Args[5])
	if err != nil {
		log.Fatalf("error reading certificate: %v", err)
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          roots,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}
}

// TLSClientConfig provides client config with updated cert and key
func TLSClientConfig() *tls.Config {
	caCertPEM, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		log.Fatalf("error reading Cert certificate: %v", ok)
	}

	cert, err := tls.LoadX509KeyPair(os.Args[6], os.Args[7])
	if err != nil {
		log.Fatalf("error reading Key certificate: %v", err)
	}
	return &tls.Config{
		RootCAs:            roots,
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}
}
