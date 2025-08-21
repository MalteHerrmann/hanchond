package noble

import (
	"errors"

	"github.com/hanchon/hanchond/playground/utils"
)

func (n *Noble) UpdateNobleGenesisFile() error {
	genesisFile, err := n.OpenGenesisFile()
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

	genesisFile, err = setAuthority(genesisFile)
	if err != nil {
		return err
	}

	genesisFile, err = setStakingDenom(genesisFile)
	if err != nil {
		return err
	}

	genesisFile, err = setWormholeState(genesisFile)
	if err != nil {
		return err
	}

	return n.SaveGenesisFile(genesisFile)
}

// addUUSDCToBank adds the denomination metadata for micro-USDC
// to the bank module's genesis state.
func addUUSDCToBank(genState map[string]any) (map[string]any, error) {
	appState, err := utils.GetAppStateFromGenesis(genState)
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
	appState, err := utils.GetAppStateFromGenesis(genState)
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
	appState, err := utils.GetAppStateFromGenesis(genState)
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

func setStakingDenom(genState map[string]any) (map[string]any, error) {
	appState, err := utils.GetAppStateFromGenesis(genState)
	if err != nil {
		return nil, err
	}

	staking, found := appState["staking"]
	if !found {
		return nil, errors.New("staking not found")
	}

	stakingGen, ok := staking.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected staking")
	}

	params, found := stakingGen["params"]
	if !found {
		return nil, errors.New("staking params not found")
	}

	stakingParams, ok := params.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected staking params")
	}

	stakingParams["bond_denom"] = "ustake"

	return genState, nil
}

func setAuthority(genState map[string]any) (map[string]any, error) {
	appState, err := utils.GetAppStateFromGenesis(genState)
	if err != nil {
		return nil, err
	}

	authority, found := appState["authority"]
	if !found {
		return nil, errors.New("authority not found")
	}

	authorityGen, ok := authority.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected authority")
	}

	// NOTE: this is a dummy account used for testing; mnemonic:
	// occur subway woman achieve deputy rapid museum point usual appear oil blue rate title claw
	// debate flag gallery level object baby winner erase carbon
	authorityGen["owner"] = "noble1zw7vatnx0vla7gzxucgypz0kfr6965akpvzw69"

	return genState, nil
}

// TODO: refactor into types and use json.Unmarshal

func setWormholeState(genState map[string]any) (map[string]any, error) {
	appState, err := utils.GetAppStateFromGenesis(genState)
	if err != nil {
		return nil, err
	}

	wormhole, found := appState["wormhole"]
	if !found {
		return nil, errors.New("wormhole not found")
	}

	wormholeGen, ok := wormhole.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected wormhole")
	}

	config, found := wormholeGen["config"]
	if !found {
		return nil, errors.New("config not found")
	}

	configGen, ok := config.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected config")
	}

	configGen["chain_id"] = 4009
	configGen["gov_chain"] = 1
	configGen["gov_address"] = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ="

	wormholeGen["guardian_sets"] = map[string]map[string]any{
		"0": {
			"addresses":       []string{"vvpCnVfNGLf4pNkaLamrSvBdD74="},
			"expiration_time": 0,
		},
	}

	return genState, nil
}
