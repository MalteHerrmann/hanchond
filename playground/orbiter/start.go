package orbiter

import (
	"fmt"
)

func (o *Orbiter) Start() (int, error) {
	cmd := fmt.Sprintf(
		"%s start --home %s --api.enable --grpc.enable >> %s 2>&1",
		o.GetVersionedBinaryPath(),
		o.HomeDir,
		o.GetLogPath(),
	)

	return o.Daemon.Start(cmd)
}
