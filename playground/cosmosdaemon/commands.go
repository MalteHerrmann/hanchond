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

func (d *Daemon) Start(options StartOptions) (int, error) {
	startCmd := d.GetVersionedBinaryPath() + " start --home " + d.HomeDir

	if d.chainInfo.IsEVMChain() {
		startCmd += " --chain-id=" + d.ChainID
	}

	startCmd += " --api.enable --grpc.enable"

	if d.chainInfo.StartFlags != "" {
		startCmd += " " + d.chainInfo.StartFlags
	}

	if options.LogLevel != "" {
		startCmd += fmt.Sprintf(" --log_level %q", options.LogLevel)
	}

	// write output to log file
	startCmd = fmt.Sprintf("%s >> %s 2>&1", startCmd, d.GetLogPath())

	return d.RunStartCmd(startCmd)
}

func (d *Daemon) RunStartCmd(startCmd string) (int, error) {
	command := exec.Command("bash", "-c", startCmd)

	// Detach the program
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err := command.Start()
	if err != nil {
		return 0, err
	}

	// repeatedly check for the child process until it's found
	var id int
	for range 10 {
		id, err = filesmanager.GetChildPID(command.Process.Pid)
		if err == nil && id != 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to find process after 10 seconds: %w", err)
	}

	return id, nil
}

// SendIBC send a simple IBC transfer using the CLI.
//
// TODO: simulate maybe? to avoid hardcoding fees.
func (d *Daemon) SendIBC(
	port, channel, recipient string,
	amount types.Coin,
	memo string,
) (string, error) {
	feeAmount := 1000
	// if the base denom is an atto unit, use a higher amount of fees to send with the transaction.
	if strings.HasPrefix(d.BaseDenom, "a") {
		feeAmount = 10000000000000000
	}

	if d.Ports == nil {
		return "", errors.New("ports are not set")
	}

	commandSlice := []string{
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
