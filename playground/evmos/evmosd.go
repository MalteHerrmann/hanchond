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
	"m/44'/60'/0'/0/0",
	types.EthAlgo,
	types.EvmosSDK,
)

type Evmos struct {
	*cosmosdaemon.Daemon
}

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
