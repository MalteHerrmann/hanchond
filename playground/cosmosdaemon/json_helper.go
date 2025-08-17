package cosmosdaemon

import (
	"encoding/json"

	"github.com/hanchon/hanchond/playground/filesmanager"
)

func readJSONFile(path string) (map[string]any, error) {
	bytes, err := filesmanager.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func saveJSONFile(data map[string]any, path string) error {
	values, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return filesmanager.SaveFile(values, path)
}
