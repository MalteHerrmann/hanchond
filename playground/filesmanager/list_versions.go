package filesmanager

import (
	"os"
)

// TODO: rename to make more general
func GetAllEvmosdVersions() ([]string, error) {
	files, err := os.Open(GetBuildsDir())
	if err != nil {
		return []string{}, err
	}
	defer files.Close()
	fileInfos, err := files.Readdir(-1)
	if err != nil {
		return []string{}, err
	}
	res := make([]string, len(fileInfos))
	for i, fileInfos := range fileInfos {
		res[i] = fileInfos.Name()
	}
	return res, nil
}
