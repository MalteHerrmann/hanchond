package evmos

import (
	"fmt"
)

func (e *Evmos) Start() (int, error) {
	logFile := e.HomeDir + "/run.log"
	cmd := fmt.Sprintf("%s start --chain-id %s --home %s --json-rpc.api eth,txpool,personal,net,debug,web3 --json-rpc.enable --api.enable --grpc.enable >> %s 2>&1",
		e.GetVersionedBinaryPath(),
		e.ChainID,
		e.HomeDir,
		logFile,
	)
	return e.Daemon.Start(cmd)
}
