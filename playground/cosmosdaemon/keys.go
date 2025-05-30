package cosmosdaemon

import (
	"fmt"
	"os/exec"
	"strings"
)

func (d *Daemon) AddValidatorKey() error {
	return d.KeyAdd(d.ValKeyName, d.ValMnemonic)
}

func (d *Daemon) KeyAdd(name string, mnemonic string) error {
	cmd := fmt.Sprintf("echo \"%s\" | %s keys add %s --recover --keyring-backend %s --home %s --key-type %s",
		mnemonic,
		d.GetVersionedBinaryPath(),
		name,
		d.KeyringBackend,
		d.HomeDir,
		d.chainInfo.GetKeyAlgo(),
	)
	command := exec.Command("bash", "-c", cmd)
	o, err := command.CombinedOutput()
	if strings.Contains(string(o), "duplicated") {
		return fmt.Errorf("duplicated address")
	}
	return err
}
