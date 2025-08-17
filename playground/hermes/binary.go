package hermes

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/types"
)

func (h *Hermes) GetHermesBinary() string {
	return filesmanager.GetHermesBinary()
}

func (h *Hermes) AddRelayerKeyIfMissing(chainID, mnemonic string, hd types.HDPath) error {
	// TODO: support account index here instead of hardwiring 0? Probably not necessary..
	hdPath := fmt.Sprintf(" --hd-path \"%s/0\" ", hd)

	logsFile := filesmanager.GetHermesPath() + "/logs_keys_" + chainID
	cmd := fmt.Sprintf(
		"echo \"%s\" | %s --config %s keys add %s --mnemonic-file /dev/stdin --chain %s >> %s 2>&1",
		mnemonic,
		h.GetHermesBinary(),
		h.GetConfigFile(),
		hdPath,
		chainID,
		logsFile,
	)

	//nolint:gosec // okay to launch subcommand here with variable
	command := exec.Command("bash", "-c", cmd)
	_, err := command.CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf(
			"%w: logs written to %s; error from logs: %s",
			err,
			logsFile,
			getErrorFromHermesLogs(logsFile),
		)
	}

	return nil
}

func (h *Hermes) CreateChannel(firstChainID, secondChainID string) error {
	logsFile := fmt.Sprintf(
		"%s/logs_channel_%s_%s",
		filesmanager.GetHermesPath(),
		firstChainID,
		secondChainID,
	)
	cmd := fmt.Sprintf(
		"%s --config %s create channel --a-chain %s --b-chain %s --a-port transfer --b-port transfer --new-client-connection --yes >> %s 2>&1",
		h.GetHermesBinary(),
		h.GetConfigFile(),
		firstChainID,
		secondChainID,
		logsFile,
	)

	//nolint:gosec // okay to launch subcommand here with variable
	command := exec.Command("bash", "-c", cmd)
	out, err := command.CombinedOutput()
	if err != nil {
		errorFromLogs := getErrorFromHermesLogs(logsFile)
		err = fmt.Errorf(
			"error %s: %s; logs written to %s; error from logs: %s",
			err.Error(),
			string(out),
			logsFile,
			errorFromLogs,
		)
	}

	return err
}

func (h *Hermes) Start() (int, error) {
	cmd := fmt.Sprintf(
		"%s --config %s start >> %s 2>&1",
		h.GetHermesBinary(),
		h.GetConfigFile(),
		filesmanager.GetHermesPath()+"/run.log",
	)

	//nolint:gosec // okay to launch subcommand here
	command := exec.Command("bash", "-c", cmd)

	// Deattach the program
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	err := command.Start()
	if err != nil {
		return 0, err
	}

	// Let hermes start
	time.Sleep(2 * time.Second)

	id, err := filesmanager.GetChildPID(command.Process.Pid)
	if err != nil {
		return 0, err
	}

	return id, nil
}
