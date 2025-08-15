package evmos

import (
	"fmt"
	"os/exec"
)

// TODO: this should be refactored to the default Cosmos daemon
func (e *Evmos) CheckBalance(wallet string) (string, error) {
	command := exec.Command( //nolint:gosec
		e.GetVersionedBinaryPath(),
		"q",
		"bank",
		"balances",
		wallet,
		"--home",
		e.HomeDir,
		"--node",
		fmt.Sprintf("http://localhost:%d", e.Ports.P26657),
	)
	out, err := command.CombinedOutput()
	return string(out), err
}
