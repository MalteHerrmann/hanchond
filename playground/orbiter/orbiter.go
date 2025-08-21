package orbiter

import (
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/types"
)

var ChainInfo = types.NewChainInfo(
	"noble",
	"simd",
	"orbiter-",
	"simd",
	"stake",
	"https://github.com/noble-assets/orbiter",
	types.CosmosHDPath,
	types.CosmosAlgo,
	"",
)

var _ cosmosdaemon.IDaemon = &Orbiter{}

type Orbiter struct {
	*cosmosdaemon.Daemon
}

func NewOrbiter(moniker, version, homeDir, chainID, keyName string, ports *types.Ports) *Orbiter {
	o := &Orbiter{
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

	o.SetCustomConfig(o.UpdateOrbiterGenesisFile)

	return o
}
