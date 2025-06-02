package curl

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

//go:embed misc/cacert.pem
var embeddedCACertData []byte

var (
	caCertPath string
	onceCACert sync.Once
	caCertErr  error
)

func getEmbeddedCACertPath() (string, error) {
	onceCACert.Do(func() {
		if len(embeddedCACertData) == 0 {
			caCertErr = fmt.Errorf("embedded CA cert data is empty, ensure cacert.pem is embedded correctly")
			return
		}

		cacheDir, err := os.UserCacheDir()
		if err != nil {
			caCertErr = fmt.Errorf("failed to get user cache dir for CA cert: %w", err)
			return
		}
		tempDir := filepath.Join(cacheDir, "go-curl-impersonate", "cacerts")
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			caCertErr = fmt.Errorf("failed to create temp dir for CA cert '%s': %w", tempDir, err)
			return
		}

		caCertPath = filepath.Join(tempDir, "embedded-cacert.pem")

		err = os.WriteFile(caCertPath, embeddedCACertData, 0644)
		if err != nil {
			caCertErr = fmt.Errorf("failed to write embedded CA cert to '%s': %w", caCertPath, err)
			return
		}
	})
	if caCertErr != nil {
		return "", caCertErr
	}

	return caCertPath, nil
}
