package utils

import "errors"

func GetAppStateFromGenesis(genState map[string]any) (map[string]any, error) {
	appState, found := genState["app_state"]
	if !found {
		return nil, errors.New("app state not found")
	}

	appStateContents, ok := appState.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected app state")
	}

	return appStateContents, nil
}
