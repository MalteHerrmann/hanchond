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
	"--json-rpc.api eth,txpool,personal,net,debug,web3 --json-rpc.enable",
)

var _ cosmosdaemon.IDaemon = &Evmos{}

type Evmos struct {
	*cosmosdaemon.Daemon
}

func NewEvmos(moniker, version, homeDir, chainID, keyName string, ports *types.Ports) *Evmos {
	e := &Evmos{
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
	e.SetCustomConfig(e.UpdateAppFile)

	return e
}
