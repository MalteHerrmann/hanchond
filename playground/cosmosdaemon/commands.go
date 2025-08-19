package cosmosdaemon

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/types"
)

func (d *Daemon) AddGenesisAccount(validatorAddr string) error {
	args := []string{
		"add-genesis-account",
		validatorAddr,
		d.ValidatorInitialSupply + d.BaseDenom,
		"--keyring-backend",
		d.KeyringBackend,
		"--home",
		d.HomeDir,
	}

	genesisSubcommand, err := d.GetGenesisSubcommand(args)
	if err != nil {
		return err
	}

	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		genesisSubcommand...,
	)

	out, err := command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(out))
	}

	return err
}

func (d *Daemon) ValidatorGenTx() error {
	args := []string{
		"gentx",
		d.ValKeyName,
		d.ValidatorInitialSupply[0:len(d.ValidatorInitialSupply)-4] + d.BaseDenom,
		"--gas-prices",
		d.BaseFee + d.BaseDenom,
		"--chain-id",
		d.ChainID,
		"--keyring-backend",
		d.KeyringBackend,
		"--home",
		d.HomeDir,
	}

	genesisSubcommand, err := d.GetGenesisSubcommand(args)
	if err != nil {
		return err
	}

	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		genesisSubcommand...,
	)
	out, err := command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(out))
	}

	return err
}

func (d *Daemon) CollectGenTxs() error {
	args := []string{
		"collect-gentxs",
		"--home",
		d.HomeDir,
	}

	genesisSubcommand, err := d.GetGenesisSubcommand(args)
	if err != nil {
		return err
	}

	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		genesisSubcommand...,
	)
	out, err := command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(out))
	}

	return err
}

func (d *Daemon) ValidateGenesis() error {
	args := []string{
		"validate-genesis",
		"--home",
		d.HomeDir,
	}

	genesisSubcommand, err := d.GetGenesisSubcommand(args)
	if err != nil {
		return err
	}

	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		genesisSubcommand...,
	)
	out, err := command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(out))
	}

	return err
}

// Returns bech32 encoded validator addresss.
func (d *Daemon) GetValidatorAddress() (string, error) {
	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		"keys",
		"show",
		"-a",
		d.ValKeyName,
		"--keyring-backend",
		d.KeyringBackend,
		"--home",
		d.HomeDir,
	)
	o, err := command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(o))

		return "", err
	}

	return strings.TrimSpace(string(o)), nil
}

func (d *Daemon) Start(startCmd string) (int, error) {
	command := exec.Command("bash", "-c", startCmd)
	// Deattach the program
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err := command.Start()
	if err != nil {
		return 0, err
	}
	time.Sleep(2 * time.Second)
	id, err := filesmanager.GetChildPID(command.Process.Pid)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SendIBC send a simple IBC transfer using the CLI.
//
// TODO: simulate maybe? to avoid hardcoding fees.
func (d *Daemon) SendIBC(port, channel, recipient string, amount types.Coin, memo string) (string, error) {
	feeAmount := 1000
	// if the base denom is an atto unit, use a higher amount of fees to send with the transaction.
	if strings.HasPrefix(d.BaseDenom, "a") {
		feeAmount = 10000000000000000
	}

	if d.Ports == nil {
		return "", errors.New("ports are not set")
	}

	commandSlice := []string{
		"tx",
		"ibc-transfer",
		"transfer",
		port,
		channel,
		recipient,
		amount.String(),
		"--keyring-backend",
		d.KeyringBackend,
		"--chain-id",
		d.ChainID,
		"--home",
		d.HomeDir,
		"--node",
		fmt.Sprintf("http://localhost:%d", d.Ports.P26657),
		"--from",
		d.ValKeyName,
		"--fees",
		fmt.Sprintf("%d%s", feeAmount, d.BaseDenom),
		"-y",
	}

	if memo != "" {
		commandSlice = append(commandSlice, "--memo", memo)
	}

	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		commandSlice...,
	)

	out, err := command.CombinedOutput()

	return string(out), err
}

func (d *Daemon) GetNodeID() (string, error) {
	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		"tendermint",
		"show-node-id",
		"--home",
		d.HomeDir,
	)
	o, err := command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(o))

		return "", err
	}

	return strings.TrimSpace(string(o)), nil
}
