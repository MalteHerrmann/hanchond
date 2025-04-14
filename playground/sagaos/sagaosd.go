package sagaos

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/filesmanager"
)

type SagaOS struct {
	*cosmosdaemon.Daemon
}

func NewSagaOS(moniker string, version string, homeDir string, chainID string, keyName string, denom string) *SagaOS {
	daemonName := version
	if !strings.Contains(version, "sagaosd") {
		daemonName = fmt.Sprintf("sagaosd%s", version)
	}
	s := &SagaOS{
		Daemon: cosmosdaemon.NewDameon(
			moniker,
			daemonName,
			homeDir,
			chainID,
			keyName,
			cosmosdaemon.EthAlgo,
			denom,
			"saga",
			cosmosdaemon.EvmosSDK,
		),
	}
	s.SetBinaryPath(filesmanager.GetDaemondPath(daemonName))
	s.SetCustomConfig(s.UpdateAppFile)
	return s
}
