package evmos

import (
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/types"
)

var ChainInfo = types.NewChainInfo(
	"evmos",
	"evmosd",
	"evmos_9000-",
	"evmosd",
	"aevmos",
	"https://github.com/evmos/evmos",
	types.EthHDPath,
	types.EthAlgo,
)

type Evmos struct {
	*cosmosdaemon.Daemon
}

// TODO: this should be able to be removed after fully refactoring to use the chain config
func NewEvmos(moniker, version, homeDir, chainID, keyName string) *Evmos {
	e := &Evmos{
		Daemon: cosmosdaemon.NewDameon(
			ChainInfo,
			moniker,
			version,
			homeDir,
			chainID,
			keyName,
		),
	}
	e.SetCustomConfig(e.UpdateAppFile)
	return e
}
