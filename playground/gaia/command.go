package gaia

import (
	"fmt"
)

func (g *Gaia) Start() (int, error) {
	logFile := g.HomeDir + "/run.log"
	cmd := fmt.Sprintf("%s start --home %s --api.enable --grpc.enable >> %s 2>&1",
		g.GetVersionedBinaryPath(),
		g.HomeDir,
		logFile,
	)
	return g.Daemon.Start(cmd)
}
