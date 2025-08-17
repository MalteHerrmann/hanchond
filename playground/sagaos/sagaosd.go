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
)

type SagaOS struct {
	*cosmosdaemon.Daemon
}

func NewSagaOS(moniker, version, homeDir, chainID, keyName string) *SagaOS {
	s := &SagaOS{
		Daemon: cosmosdaemon.NewDameon(
			ChainInfo,
			moniker,
			version,
			homeDir,
			chainID,
			keyName,
		),
	}
	s.SetCustomConfig(s.UpdateAppFile)

	return s
}
