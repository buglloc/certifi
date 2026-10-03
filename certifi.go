package certifi

import (
	"crypto/x509"
)

// NewCertPool returns a new pool containing the builtin certificates.
// It panics if the embedded bundle is invalid or contains no certificates.
func NewCertPool() *x509.CertPool {
	out := x509.NewCertPool()
	for _, cert := range certs() {
		out.AddCert(cert)
	}

	return out
}
