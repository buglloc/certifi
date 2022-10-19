package certifi_test

import (
	"testing"

	"github.com/buglloc/certifi"
)

func TestParseCerts(t *testing.T) {
	certs, err := certifi.ParseCerts()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(certs) < 100 {
		t.Errorf("unexpected certs cound %d < 100", len(certs))
	}
}
