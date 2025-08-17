package gaia

import (
	"fmt"
)

func (g *Gaia) Start() (int, error) {
	cmd := fmt.Sprintf(
		"%s start --home %s --api.enable --grpc.enable >> %s 2>&1",
		g.GetVersionedBinaryPath(),
		g.HomeDir,
		g.GetLogPath(),
	)

	return g.Daemon.Start(cmd)
}
