package certifi

import (
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"sync"
)

//go:generate go run ./cmd/download --out cacert.pem

//go:embed cacert.pem
var certsBytes []byte

var (
	parseCertsOnce sync.Once
	parsedCerts    []*x509.Certificate
	parseCertsErr  error
)

func certs() []*x509.Certificate {
	parseCertsOnce.Do(func() {
		parsedCerts, parseCertsErr = ParseCertificates(certsBytes)
		if parseCertsErr == nil && len(parsedCerts) == 0 {
			parseCertsErr = errors.New("bundle contains no certificates")
		}
	})

	// Panic outside Once so callers that recover cannot subsequently get an empty pool.
	if parseCertsErr != nil {
		panic(fmt.Errorf("certifi: invalid builtin certificate bundle: %w", parseCertsErr))
	}

	return parsedCerts
}

// ParseCertificates parses CERTIFICATE PEM blocks from in, preserving their order.
// It ignores non-PEM text and skips blocks with other types or PEM headers.
//
// Empty input returns nil and io.ErrUnexpectedEOF. Invalid certificate DER
// returns nil and a wrapped parsing error, discarding any certificates already parsed.
// Nonempty input with no certificate blocks returns nil, nil.
func ParseCertificates(in []byte) ([]*x509.Certificate, error) {
	if len(in) == 0 {
		return nil, io.ErrUnexpectedEOF
	}

	var out []*x509.Certificate
	rest := in
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
