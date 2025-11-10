package evmos

import (
	"fmt"
	"os/exec"
)

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
		"10000000000000000"+e.BaseDenom,
		"-y",
	)

	out, err := command.CombinedOutput()

	return string(out), err
}
