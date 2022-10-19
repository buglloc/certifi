package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fatalf(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "certifi-download: "+msg+"\n", a...)
	os.Exit(1)
}

func main() {
	var out string
	flag.StringVar(&out, "out", "cacert.pem", "output filename")
	flag.Parse()

	outF, err := os.Create(out)
	if err != nil {
		fatalf("unable to create output file: %v", err)
	}
	_, _ = outF.WriteString(fmt.Sprintf("# Generated at: %s\n\n", time.Now().UTC().Format(time.RFC3339)))

	defer func() {
		if err := outF.Close(); err != nil {
			fatalf("unable to close output file: %v", err)
		}
	}()

	resp, err := http.Get("https://mkcert.org/generate")
	if err != nil {
		fatalf("can't make HTTP request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		fatalf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(outF, resp.Body)
	if err != nil {
		fatalf("download failed: %s", resp.Status)
	}
}
