package certifi

import (
	"crypto/x509"
)

func NewCertPool() *x509.CertPool {
	out := x509.NewCertPool()
	for _, cert := range certs() {
		out.AddCert(cert)
	}

	return out
}
