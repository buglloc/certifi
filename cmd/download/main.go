package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/buglloc/certifi"
)

func fatalf(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "certifi-download: "+msg+"\n", a...)

	os.Exit(1)
}

func main() {
	var out string
	flag.StringVar(&out, "out", "cacert.pem", "output filename")
	flag.Parse()

	if err := download(out, "https://mkcert.org/generate"); err != nil {
		fatalf("%v", err)
	}
}

func download(out, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	outF, err := os.CreateTemp(filepath.Dir(out), "."+filepath.Base(out)+"-*")
	if err != nil {
		return fmt.Errorf("unable to create temporary output file: %w", err)
	}
	defer os.Remove(outF.Name())
	defer outF.Close()

	if _, err := fmt.Fprintf(outF, "# Generated at: %s\n\n", time.Now().UTC().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("unable to write output header: %w", err)
	}

	if _, err := io.Copy(outF, resp.Body); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if err := outF.Close(); err != nil {
		return fmt.Errorf("unable to close output file: %w", err)
	}

	bundle, err := os.ReadFile(outF.Name())
	if err != nil {
		return fmt.Errorf("unable to read downloaded bundle: %w", err)
	}

	certs, err := certifi.ParseCertificates(bundle)
	if err != nil {
		return fmt.Errorf("invalid downloaded bundle: %w", err)
	}

	if len(certs) == 0 {
		return fmt.Errorf("downloaded bundle contains no certificates")
	}

	mode := os.FileMode(0644)

	info, err := os.Stat(out)
	if err == nil {
		mode = info.Mode().Perm()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("unable to stat output file: %w", err)
	}

	if err := os.Chmod(outF.Name(), mode); err != nil {
		return fmt.Errorf("unable to set output permissions: %w", err)
	}

	if err := os.Rename(outF.Name(), out); err != nil {
		return fmt.Errorf("unable to replace output file: %w", err)
	}

	return nil
}
