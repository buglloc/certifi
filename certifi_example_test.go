package certifi_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/buglloc/certifi"
)

func ExampleNewCertPool_http() {
	httpc := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certifi.NewCertPool(),
			},
		},
	}

	resp, err := httpc.Get("https://google.com/")
	if err != nil {
		panic(fmt.Sprintf("http failed: %v\n", err))
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	fmt.Printf("successful, status code: %d\n", resp.StatusCode)
}
