package cosmosdaemon

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/converter"
	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/txbuilder"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/types"
)

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

	Ports *Ports

	CustomConfig func() error
}

func NewDameon(
	chainInfo types.ChainInfo,
	moniker string,
	version string,
	homeDir string,
	chainID string,
	keyName string,
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

		Ports: nil,
	}
}

func (d *Daemon) GetVersionedBinaryPath() string {
	return filesmanager.GetDaemondPathWithVersion(d.chainInfo, d.Version)
}

// This is used to change the config files that are specific to a client
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
