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
)

type Orbiter struct {
	*cosmosdaemon.Daemon
}

func NewOrbiter(moniker, version, homeDir, chainID, keyName string) *Orbiter {
	o := &Orbiter{
		Daemon: cosmosdaemon.NewDameon(
			ChainInfo,
			moniker,
			version,
			homeDir,
			chainID,
			keyName,
		),
	}
	o.SetCustomConfig(o.UpdateOrbiterGenesisFile)

	return o
}
