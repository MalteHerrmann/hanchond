package types

import (
	"encoding/json"
	"strings"
)

type ChainInfo struct {
	AccountPrefix string `json:"account_prefix"`
	BinaryName    string `json:"binary_name"`
	ChainIDBase   string `json:"chain_id_base"`
	// TODO: is this even used somewhere? maybe remove or repurpose
	ClientName string        `json:"client_name"`
	Denom      string        `json:"denom"`
	HDPath     HDPath        `json:"hd_path,omitempty"`
	KeyAlgo    SignatureAlgo `json:"key_algo"`
	RepoURL    string        `json:"repo_url"`
	StartFlags string        `json:"start_flags"`
}

func NewChainInfo(
	accountPrefix string,
	binaryName string,
	chainIDBase string,
	clientName string,
	denom string,
	repoURL string,
	hdPath HDPath,
	keyAlgo SignatureAlgo,
	startFlags string,
) ChainInfo {
	return ChainInfo{
		AccountPrefix: accountPrefix,
		BinaryName:    binaryName,
		ChainIDBase:   chainIDBase,
		ClientName:    clientName,
		Denom:         denom,
		HDPath:        hdPath,
		RepoURL:       repoURL,
		KeyAlgo:       keyAlgo,
		StartFlags:    startFlags,
	}
}

func MustParseChainInfo(input string) ChainInfo {
	ci, err := ParseChainInfo(input)
	if err != nil {
		panic("could not parse chain info: " + err.Error())
	}

	return ci
}

func (ci ChainInfo) GetAccountPrefix() string  { return ci.AccountPrefix }
func (ci ChainInfo) GetBinaryName() string     { return ci.BinaryName }
func (ci ChainInfo) GetChainIDBase() string    { return ci.ChainIDBase }
func (ci ChainInfo) GetClientName() string     { return ci.ClientName }
func (ci ChainInfo) GetDenom() string          { return ci.Denom }
func (ci ChainInfo) GetHDPath() HDPath         { return ci.HDPath }
func (ci ChainInfo) GetKeyAlgo() SignatureAlgo { return ci.KeyAlgo }
func (ci ChainInfo) GetRepoURL() string        { return ci.RepoURL }
func (ci ChainInfo) GetStartFlags() string     { return ci.StartFlags }

func (ci ChainInfo) IsEVMChain() bool {
	return ci.KeyAlgo == EthAlgo
}

func (ci ChainInfo) GetPrefixedDaemonName(version string) string {
	daemonName := version
	if !strings.Contains(version, ci.BinaryName) {
		daemonName = ci.BinaryName + version
	}

	return daemonName
}

func (ci ChainInfo) GetVersionedBinaryName(version string) string {
	return ci.BinaryName + "_" + version
}

func (ci ChainInfo) MustMarshal() []byte {
	bz, err := json.Marshal(&ci)
	if err != nil {
		panic(err)
	}

	return bz
}

func ParseChainInfo(input string) (ci ChainInfo, err error) {
	err = json.Unmarshal([]byte(input), &ci)

	return
}

type SignatureAlgo string

const (
	EthAlgo    SignatureAlgo = "eth_secp256k1"
	CosmosAlgo SignatureAlgo = "secp256k1"
)

type HDPath string

const (
	// CosmosHDPath as per:
	// https://github.com/confio/cosmos-hd-key-derivation-spec?tab=readme-ov-file#the-cosmos-hub-path
	CosmosHDPath HDPath = "m/44'/118'/0'/0"
	// EthHDPath as per e.g.: https://docs.ethers.org/v5/api/utils/hdnode/
	EthHDPath HDPath = "m/44'/60'/0'/0"
)
