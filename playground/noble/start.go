package noble

import "fmt"

func (n *Noble) Start() (int, error) {
	return n.Daemon.Start(
		fmt.Sprintf(
			"%s start --home %s --api.enable --grpc.enable >> %s 2>&1",
			n.GetVersionedBinaryPath(),
			n.HomeDir,
			n.GetLogPath(),
		),
	)
}
