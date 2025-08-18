package orbiter

import (
	"errors"
)

// UpdateOrbiterGenesisFile extends the default Cosmos daemon's method
// to set the required items for the orbiter simapp.
func (o *Orbiter) UpdateOrbiterGenesisFile() error {
	genesisFile, err := o.OpenGenesisFile()
	if err != nil {
		return err
	}

	genesisFile, err = addUUSDCToBank(genesisFile)
	if err != nil {
		return err
	}

	genesisFile, err = setFTFState(genesisFile)
	if err != nil {
		return err
	}

	genesisFile, err = setCCTPState(genesisFile)
	if err != nil {
		return err
	}

	return o.SaveGenesisFile(genesisFile)
}

// addUUSDCToBank adds the denomination metadata for micro-USDC
// to the bank module's genesis state.
func addUUSDCToBank(genState map[string]any) (map[string]any, error) {
	appState, err := getAppStateFromGenesis(genState)
	if err != nil {
		return nil, err
	}

	bankState, found := appState["bank"]
	if !found {
		return nil, errors.New("bank app state not found")
	}

	bankContents, ok := bankState.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected bank state")
	}

	bankContents["denom_metadata"] = []map[string]any{{
		"description": "Circle USD Coin",
		"denom_units": []map[string]any{
			{"denom": "uusdc", "exponent": 0, "aliases": []string{"microusdc"}},
			{"denom": "usdc", "exponent": 6},
		},
		"base":    "uusdc",
		"display": "usdc",
		"name":    "Circle USD Coin",
		"symbol":  "USDC",
	}}

	return genState, nil
}

// setFTFState updates the genesis state for CCTP's
// fiat token factory, which is required for the orbiter simapp
// to run.
func setFTFState(genState map[string]any) (map[string]any, error) {
	appState, err := getAppStateFromGenesis(genState)
	if err != nil {
		return nil, err
	}

	ftfState, found := appState["fiat-tokenfactory"]
	if !found {
		return nil, errors.New("fiat-tokenfactory not found")
	}

	ftfGenesis, ok := ftfState.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected fiat-tokenfactory")
	}

	ftfGenesis["mintingDenom"] = map[string]any{"denom": "uusdc"}
	ftfGenesis["paused"] = map[string]any{"paused": false}
	ftfGenesis["mintersList"] = []map[string]any{
		{
			"address":   "noble12l2w4ugfz4m6dd73yysz477jszqnfughxvkss5",
			"allowance": map[string]any{"denom": "uusdc", "amount": "1000000000000"},
		},
	}

	return genState, nil
}

func setCCTPState(genState map[string]any) (map[string]any, error) {
	appState, err := getAppStateFromGenesis(genState)
	if err != nil {
		return nil, err
	}

	cctp, found := appState["cctp"]
	if !found {
		return nil, errors.New("cctp not found")
	}

	cctpGenesis, ok := cctp.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected cctp")
	}

	cctpGenesis["per_message_burn_limit_list"] = []map[string]any{
		{"denom": "uusdc", "amount": "1000000000000"},
	}
	cctpGenesis["token_messenger_list"] = []map[string]any{
		{
			"domain_id": 0,
			"address":   "AAAAAAAAAAAAAAAAvT+oG1i6kqghNgOLJa3scGavMVU=",
		},
	}

	return genState, nil
}

func getAppStateFromGenesis(genState map[string]any) (map[string]any, error) {
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
