package gaia

import (
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/types"
)

var ChainInfo = types.NewChainInfo(
	"cosmos",
	"gaiad",
	"cosmoshub-",
	"gaiad",
	"icsstake",
	"https://github.com/cosmos/gaia",
	types.CosmosAlgo,
	types.GaiaSDK,
)

type Gaia struct {
	*cosmosdaemon.Daemon
}

func NewGaia(moniker, homeDir, chainID, keyName string) *Gaia {
	g := &Gaia{
		Daemon: cosmosdaemon.NewDameon(
			ChainInfo,
			moniker,
			// TODO: enable using different versions in Gaia?
			"gaia",
			homeDir,
			chainID,
			keyName,
		),
	}
	g.SetCustomConfig(g.UpdateGenesisFile)
	return g
}
