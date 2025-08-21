package sagaos

import (
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/types"
)

var ChainInfo = types.NewChainInfo(
	"saga",
	"sagaosd",
	"sagaos_1234-",
	// TODO: is this even used? check
	"sagaosd",
	"saga",
	"https://github.com/sagaxyz/sagaos",
	types.CosmosHDPath,
	types.EthAlgo,
	" --json-rpc.enable true --json-rpc.api eth,txpool,personal,net,debug,web3",
)

var _ cosmosdaemon.IDaemon = &SagaOS{}

type SagaOS struct {
	*cosmosdaemon.Daemon
}

func NewSagaOS(moniker, version, homeDir, chainID, keyName string, ports *types.Ports) *SagaOS {
	daemon := cosmosdaemon.NewDaemon(
		ChainInfo,
		moniker,
		version,
		homeDir,
		chainID,
		keyName,
		ports,
	)

	s := &SagaOS{
		Daemon: daemon,
	}
	s.SetCustomConfig(s.UpdateAppFile)

	return s
}
