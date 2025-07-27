package certifi

import "crypto/x509"

func ParseBuiltinCerts() ([]*x509.Certificate, error) {
	return ParseCertificates(certsBytes)
}
