package gaia

func (g *Gaia) UpdateGenesisFile() error {
	// Gaia extra config
	genesis, err := g.OpenGenesisFile()
	if err != nil {
		return err
	}

	g.setUnbondingTime(genesis)
	// Maybe we need to update the `genesis_time` but I am not sure why

	return g.SaveGenesisFile(genesis)
}

func (g *Gaia) setUnbondingTime(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app state")
	}

	if v, ok := appState["staking"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					// Base Denom
					if _, ok := v["unbonding_time"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["staking"].(map[string]any)["params"].(map[string]any)["unbonding_time"] = "1814400s"
					}
				}
			}
		}
	}
}
