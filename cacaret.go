package certifi

import (
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"fmt"
	"sync"
)

//go:generate go run ./cmd/download --out cacert.pem

//go:embed cacert.pem
var certsBytes []byte

var (
	parseCertsOnce sync.Once
	parsedCerts    []*x509.Certificate
)

func certs() []*x509.Certificate {
	parseCertsOnce.Do(func() {
		parsedCerts, _ = parseCerts()
	})
	return parsedCerts
}

func parseCerts() ([]*x509.Certificate, error) {
	var out []*x509.Certificate
	rest := certsBytes
	for {
		if len(rest) == 0 {
			break
		}

		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse certificate from block %q: %w", string(block.Bytes), err)
		}

		out = append(out, cert)
	}

	return out, nil
}
