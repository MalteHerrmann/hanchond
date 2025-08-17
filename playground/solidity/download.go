package solidity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hanchon/hanchond/playground/filesmanager"
)

type build struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}

type versionList struct {
	Builds []build `json:"builds"`
}

func DownloadSolcBinary(isDarwin bool, version string) error {
	baseURL := "https://binaries.soliditylang.org/macosx-amd64/"
	if !isDarwin {
		baseURL = strings.Replace(baseURL, "macosx", "linux", 1)
	}
	list, err := http.Get(baseURL + "list.json")
	if err != nil {
		return err
	}
	if list.StatusCode != 200 {
		return fmt.Errorf("status code not 200")
	}

	listContent, err := io.ReadAll(list.Body)
	if err != nil {
		return err
	}
	defer list.Body.Close() //nolint:errcheck

	var v versionList
	if err := json.Unmarshal(listContent, &v); err != nil {
		return err
	}

	binaryURL := ""
	for _, v := range v.Builds {
		if v.Version == version {
			binaryURL = (baseURL + v.Path)

			break
		}
	}
	if binaryURL == "" {
		return fmt.Errorf("solidity version not found")
	}

	filePathInDisk := filesmanager.GetSolcPath(version)

	//nolint:gosec // okay to create file here
	file, err := os.Create(filePathInDisk)
	if err != nil {
		return fmt.Errorf("could not create file:%s", err.Error())
	}
	defer file.Close() //nolint:errcheck

	resp, err := http.Get(binaryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error requesting the binary file: %d", resp.StatusCode)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("could not save the binary file:%s", err.Error())
	}

	// Executable
	info, err := os.Stat(filePathInDisk)
	if err != nil {
		return err
	}

	// Add execute permissions
	if err := os.Chmod(filePathInDisk, info.Mode()|0o111); err != nil {
		return err
	}

	return nil
}
