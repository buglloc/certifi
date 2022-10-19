package certifi

import "crypto/x509"

func ParseCerts() ([]*x509.Certificate, error) {
	return parseCerts()
}
