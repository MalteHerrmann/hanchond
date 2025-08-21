package noble

import (
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/types"
)

var ChainInfo = types.NewChainInfo(
	"noble",
	"nobled",
	"duke-",
	"nobled",
	"ustake",
	"https://github.com/noble-assets/noble",
	types.CosmosHDPath,
	types.CosmosAlgo,
	"",
)

var _ cosmosdaemon.IDaemon = &Noble{}

type Noble struct {
	*cosmosdaemon.Daemon
}

func NewNoble(
	moniker, version, homeDir, chainID, keyName string,
	ports *types.Ports,
) *Noble {
	n := &Noble{
		Daemon: cosmosdaemon.NewDaemon(
			ChainInfo,
			moniker,
			version,
			homeDir,
			chainID,
			keyName,
			ports,
		),
	}

	n.SetCustomConfig(n.UpdateNobleGenesisFile)

	return n
}
