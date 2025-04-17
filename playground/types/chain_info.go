package types

import (
	"encoding/json"
	"strings"
)

type ChainInfo struct {
	AccountPrefix string `json:"account_prefix"`
	BinaryName    string `json:"binary_name"`
	ChainIDBase   string `json:"chain_id_base"`
	// TODO: is this even used somewhere? maybe remove
	ClientName string        `json:"client_name"`
	Denom      string        `json:"denom"`
	HDPath     string        `json:"hd_path,omitempty"`
	KeyAlgo    SignatureAlgo `json:"key_algo"`
	SdkVersion SDKVersion    `json:"sdk_version"`
	RepoURL    string        `json:"repo_url"`
}

func (ci ChainInfo) GetAccountPrefix() string  { return ci.AccountPrefix }
func (ci ChainInfo) GetBinaryName() string     { return ci.BinaryName }
func (ci ChainInfo) GetChainIDBase() string    { return ci.ChainIDBase }
func (ci ChainInfo) GetClientName() string     { return ci.ClientName }
func (ci ChainInfo) GetDenom() string          { return ci.Denom }
func (ci ChainInfo) GetHDPath() string         { return ci.HDPath }
func (ci ChainInfo) GetKeyAlgo() SignatureAlgo { return ci.KeyAlgo }
func (ci ChainInfo) GetRepoURL() string        { return ci.RepoURL }
func (ci ChainInfo) GetSDKVersion() SDKVersion { return ci.SdkVersion }

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

func MustParseChainInfo(input string) ChainInfo {
	ci, err := ParseChainInfo(input)
	if err != nil {
		panic("could not parse chain info: " + err.Error())
	}
	return ci
}

func NewChainInfo(
	accountPrefix string,
	binaryName string,
	chainIDBase string,
	clientName string,
	denom string,
	repoURL string,
	hdPath string,
	keyAlgo SignatureAlgo,
	sdkVersion SDKVersion,
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
		SdkVersion:    sdkVersion,
	}
}

type SignatureAlgo string

const (
	EthAlgo    SignatureAlgo = "eth_secp256k1"
	CosmosAlgo SignatureAlgo = "secp256k1"
)

type SDKVersion string

const (
	// NOTE: there are some differences in the namespace while interacting with the CLI, like the genesis namespace
	GaiaSDK  SDKVersion = "gaiaSDK"
	EvmosSDK SDKVersion = "evmosSDK"
)
