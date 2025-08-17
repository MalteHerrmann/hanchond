package gaia

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

func GetGaiadBinary(isDarwin bool, version string) error {
	arch := runtime.GOARCH
	if arch != "arm64" {
		arch = "amd64"
	}
	systemOS := "darwin"
	if !isDarwin {
		systemOS = "linux"
	}

	url := fmt.Sprintf(
		"https://github.com/cosmos/gaia/releases/download/%s/gaiad-%s-%s-%s",
		version,
		version,
		systemOS,
		arch,
	)

	path := ChainInfo.GetVersionedBinaryName(version)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download Gaia: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download gaiad binary: status code %d", resp.StatusCode)
	}

	//nolint:gosec // file operation is fine here
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close() //nolint:errcheck

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save gaiad binary: %w", err)
	}

	err = os.Chmod(path, 0o600)
	if err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}
