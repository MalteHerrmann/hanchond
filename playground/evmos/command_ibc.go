package evmos

import (
	"fmt"
	"os/exec"
)

// TODO: this can be refactored to the default Cosmos daemon as we assume all chains are IBC-enabled for now
func (e *Evmos) SendIBC(port, channel, receiver, amount string) (string, error) {
	command := exec.Command( //nolint:gosec
		e.GetVersionedBinaryPath(),
		"tx",
		"ibc-transfer",
		"transfer",
		port,
		channel,
		receiver,
		amount,
		"--keyring-backend",
		e.KeyringBackend,
		"--home",
		e.HomeDir,
		"--node",
		fmt.Sprintf("http://localhost:%d", e.Ports.P26657),
		"--from",
		e.ValKeyName,
		"--fees",
		fmt.Sprintf("10000000000000000%s", e.BaseDenom),
		"-y",
	)

	out, err := command.CombinedOutput()
	return string(out), err
}
