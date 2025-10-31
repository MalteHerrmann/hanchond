package cosmosdaemon

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/converter"
	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/txbuilder"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/types"
)

// IDaemon defines the interface that all Cosmos SDK-based binaries
// that are configured in this tool should fulfill.
type IDaemon interface {
	Balance(address string) (string, error)
	Start(options StartOptions) (pid int, err error)
	SendIBC(port, channel, recipient string, amount types.Coin, memo string) (string, error)
	Tx(hash string) (string, error)
}

type Daemon struct {
	chainInfo   types.ChainInfo
	ValKeyName  string
	ValMnemonic string
	ValWallet   string
	Version     string

	KeyringBackend string
	HomeDir        string

	ChainID string
	Moniker string

	BaseDenom string
	GasLimit  string
	BaseFee   string

	ValidatorInitialSupply string

	Ports *types.Ports

	CustomConfig func() error
}

func NewDaemon(
	chainInfo types.ChainInfo,
	moniker string,
	version string,
	homeDir string,
	chainID string,
	keyName string,
	ports *types.Ports,
) *Daemon {
	mnemonic, _ := txbuilder.NewMnemonic()
	wallet := ""
	if chainInfo.IsEVMChain() {
		_, temp, _ := txbuilder.WalletFromMnemonic(mnemonic)
		wallet, _ = converter.HexToBech32(temp.Address.Hex(), chainInfo.GetAccountPrefix())
	} else {
		wallet, _ = txbuilder.MnemonicToCosmosAddress(mnemonic, chainInfo.GetAccountPrefix())
	}

	return &Daemon{
		chainInfo:   chainInfo,
		ValKeyName:  keyName,
		ValMnemonic: mnemonic,
		ValWallet:   wallet,
		// TODO: add validity check for version?
		Version: version,

		KeyringBackend: "test",
		HomeDir:        homeDir,

		ChainID: chainID,
		Moniker: moniker,

		BaseDenom: chainInfo.GetDenom(),

		ValidatorInitialSupply: "100000000000000000000000000",

		// Maybe move this to just evmos
		GasLimit: "1000000000",
		BaseFee:  "1000000000",

		Ports: ports,
	}
}

func (d *Daemon) GetChainInfo() types.ChainInfo {
	return d.chainInfo
}

func (d *Daemon) GetVersionedBinaryPath() string {
	return filesmanager.GetDaemondPathWithVersion(d.chainInfo, d.Version)
}

// SetCustomConfig is used to change the config files that are specific to a client.
func (d *Daemon) SetCustomConfig(configurator func() error) {
	d.CustomConfig = configurator
}

func (d *Daemon) ExecuteCustomConfig() error {
	if d.CustomConfig == nil {
		return nil
	}

	return d.CustomConfig()
}

func (d *Daemon) SetValidatorWallet(mnemonic, wallet string) {
	d.ValMnemonic = mnemonic
	d.ValWallet = wallet
}

func (d *Daemon) Tx(hash string) (string, error) {
	return utils.ExecCommand(
		d.GetVersionedBinaryPath(),
		"q",
		"tx",
		"--type=hash",
		hash,
		"--node",
		fmt.Sprintf("http://localhost:%d", d.Ports.P26657),
	)
}

func (d *Daemon) NewRequester() *requester.Client {
	return requester.NewClient().
		WithUnsecureWeb3Endpoint(fmt.Sprintf("http://localhost:%d", d.Ports.P8545)).
		WithUnsecureRestEndpoint(fmt.Sprintf("http://localhost:%d", d.Ports.P1317)).
		WithUnsecureTendermintEndpoint(fmt.Sprintf("http://localhost:%d", d.Ports.P26657))
}

func (d *Daemon) NewTxBuilder(gasLimit uint64) *txbuilder.TxBuilder {
	return txbuilder.NexTxBuilder(
		map[string]txbuilder.Contract{},
		d.ValMnemonic,
		map[string]uint64{},
		gasLimit,
		d.NewRequester(),
	)
}
