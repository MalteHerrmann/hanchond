package hermes

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/types"
)

func LocalEndpoint(port int64) string {
	return fmt.Sprintf("http://127.0.0.1:%d", port)
}

func (h *Hermes) GetConfigFile() string {
	// If the dir already existed it will return error, but that is fine
	_ = filesmanager.CreateHermesFolder()

	return filesmanager.GetHermesPath() + "/config.toml"
}

func (h *Hermes) AddCosmosChain(chainInfo types.ChainInfo, chainID, p26657, p9090, keyname, mnemonic string) error {
	configFile, err := filesmanager.ReadFile(h.GetConfigFile())
	if err != nil {
		return err
	}

	configFileString := string(configFile)
	// If the chain was already included in the config file, do nothing
	//
	// TODO: Check if ports need updating
	if !strings.Contains(configFileString, chainID) {
		configFileString += getCosmosChainConfig(chainInfo, chainID, p26657, p9090, keyname)

		err = filesmanager.SaveFile([]byte(configFileString), h.GetConfigFile())
		if err != nil {
			panic(err)
		}

		return nil
	}

	err = h.AddRelayerKeyIfMissing(chainID, mnemonic, chainInfo.GetHDPath())
	if err != nil {
		// NOTE: if we had a problem adding the relayer key we are overwriting the changes made to the relayer config
		_ = filesmanager.SaveFile(configFile, h.GetConfigFile())

		panic(err)
	}

	return nil
}

func (h *Hermes) AddEVMChain(chainInfo types.ChainInfo, chainID, p26657, p9090, keyname, mnemonic string) error {
	configFile, err := filesmanager.ReadFile(h.GetConfigFile())
	if err != nil {
		return err
	}

	configFileString := string(configFile)

	// We only add the chain to the configuration if it had not been added already
	//
	// TODO: Maybe check if updates to ports are required? currently it's only either add fully or nothing at all
	if !strings.Contains(configFileString, chainID) {
		configFileString += getEVMChainConfig(chainInfo, chainID, p26657, p9090, keyname)

		err = filesmanager.SaveFile([]byte(configFileString), h.GetConfigFile())
		if err != nil {
			panic(err)
		}
	}

	err = h.AddRelayerKeyIfMissing(chainID, mnemonic, chainInfo.GetHDPath())
	if err != nil {
		// NOTE: if we had a problem adding the relayer key we are overwriting the changes made to the relayer config
		_ = filesmanager.SaveFile(configFile, h.GetConfigFile())

		panic(err)
	}

	return nil
}

func (h *Hermes) initHermesConfig() {
	// Init the file only if it does not exist
	if filesmanager.DoesFileExist(h.GetConfigFile()) {
		return
	}

	basicConfig := `
[global]
log_level = 'trace'

[mode]

[mode.clients]
enabled = true
refresh = true

[mode.connections]
enabled = false

[mode.channels]
enabled = false

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true

[rest]
enabled = false
host = '127.0.0.1'
port = 3000

[telemetry]
enabled = false
host = '127.0.0.1'
port = 3001
`

	err := filesmanager.SaveFile([]byte(basicConfig), h.GetConfigFile())
	if err != nil {
		panic(err)
	}
}

func getCosmosChainConfig(chainInfo types.ChainInfo, chainID, p26657, p9090, keyname string) string {
	return fmt.Sprintf(`
[[chains]]
id = '%s'
rpc_addr = '%s'
grpc_addr = '%s'
event_source = { mode = 'pull', interval = '1s' }
rpc_timeout = '10s'
account_prefix = '%s'
key_store_folder = '%s'
key_name = '%s'
store_prefix = 'ibc'
default_gas = 5000000
max_gas = 10000000
gas_price = { price = 100, denom = '%s' }
gas_multiplier = 5
max_msg_num = 20
max_tx_size = 209715
clock_drift = '20s'
max_block_time = '10s'
trust_threshold = { numerator = '1', denominator = '3' }
`, chainID, p26657, p9090, chainInfo.GetAccountPrefix(), filesmanager.GetHermesPath()+"/keyring"+chainID, keyname, chainInfo.GetDenom()) //nolint:lll
}

func getEVMChainConfig(chainInfo types.ChainInfo, chainID, p26657, p9090, keyname string) string {
	return fmt.Sprintf(`
[[chains]]
id = '%s'
rpc_addr = '%s'
grpc_addr = '%s'
event_source = { mode = 'pull', interval = '1s' }
rpc_timeout = '10s'
key_name = '%s'
key_store_folder = '%s'
store_prefix = 'ibc'
default_gas = 100000
max_gas = 3000000
clock_drift = '15s'
max_block_time = '10s'
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
account_prefix = '%s'
gas_price = { price = 800000000, denom = '%s' }
address_type = { derivation = 'ethermint', proto_type = { pk_type = '/ethermint.crypto.v1.ethsecp256k1.PubKey' } }
`, chainID, p26657, p9090, keyname, filesmanager.GetHermesPath()+"/keyring"+chainID, chainInfo.GetAccountPrefix(), chainInfo.GetDenom()) //nolint:lll
}
