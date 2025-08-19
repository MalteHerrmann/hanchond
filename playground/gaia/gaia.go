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
	types.CosmosHDPath,
	types.CosmosAlgo,
)

var _ cosmosdaemon.IDaemon = &Gaia{}

type Gaia struct {
	*cosmosdaemon.Daemon
}

func NewGaia(moniker, homeDir, chainID, keyName string, ports *types.Ports) *Gaia {
	g := &Gaia{
		Daemon: cosmosdaemon.NewDaemon(
			ChainInfo,
			moniker,
			// TODO: enable using different versions in Gaia?
			"gaia",
			homeDir,
			chainID,
			keyName,
			ports,
		),
	}
	g.SetCustomConfig(g.UpdateGenesisFile)

	return g
}
