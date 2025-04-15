package types

import (
	"strings"
)

type ChainInfo struct {
	accountPrefix string
	binaryName    string
	chainIDBase   string
	// TODO: is this even used somewhere? maybe remove
	clientName string
	denom      string
	keyAlgo    SignatureAlgo
	sdkVersion SDKVersion
	repoURL    string
}

func (ci ChainInfo) GetAccountPrefix() string  { return ci.accountPrefix }
func (ci ChainInfo) GetBinaryName() string     { return ci.binaryName }
func (ci ChainInfo) GetChainIDBase() string    { return ci.chainIDBase }
func (ci ChainInfo) GetClientName() string     { return ci.clientName }
func (ci ChainInfo) GetDenom() string          { return ci.denom }
func (ci ChainInfo) GetKeyAlgo() SignatureAlgo { return ci.keyAlgo }
func (ci ChainInfo) GetRepoURL() string        { return ci.repoURL }
func (ci ChainInfo) GetSDKVersion() SDKVersion { return ci.sdkVersion }

func (ci ChainInfo) GetPrefixedDaemonName(version string) string {
	daemonName := version
	if !strings.Contains(version, ci.binaryName) {
		daemonName = ci.binaryName + version
	}

	return daemonName
}

func (ci ChainInfo) IsEVMChain() bool {
	return ci.keyAlgo == EthAlgo
}

func NewChainInfo(
	accountPrefix string,
	binaryName string,
	chainIDBase string,
	clientName string,
	denom string,
	repoURL string,
	keyAlgo SignatureAlgo,
	sdkVersion SDKVersion,
) ChainInfo {
	return ChainInfo{
		accountPrefix: accountPrefix,
		binaryName:    binaryName,
		chainIDBase:   chainIDBase,
		clientName:    clientName,
		denom:         denom,
		repoURL:       repoURL,
		keyAlgo:       keyAlgo,
		sdkVersion:    sdkVersion,
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
