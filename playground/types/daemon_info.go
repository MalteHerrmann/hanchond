package types

// DaemonInfo contains all relevant information that is required to
// instantiate Cosmos SDK-based daemons for CLI commands.
//
// TODO: check if this can be moved into / generated from SQL stuff maybe?
type DaemonInfo struct {
	moniker          string
	version          string
	configFolder     string
	chainID          string
	validatorKeyName string
	chainInfo        ChainInfo
	ports            Ports
}

// NewDaemonInfo returns a new DaemonInfo.
func NewDaemonInfo(
	moniker string,
	version string,
	configFolder string,
	chainID string,
	validatorKeyName string,
	chainInfo ChainInfo,
	ports Ports,
) DaemonInfo {
	return DaemonInfo{
		moniker:          moniker,
		version:          version,
		configFolder:     configFolder,
		chainID:          chainID,
		validatorKeyName: validatorKeyName,
		chainInfo:        chainInfo,
		ports:            ports,
	}
}

func (d DaemonInfo) GetChainID() string          { return d.chainID }
func (d DaemonInfo) GetChainInfo() ChainInfo     { return d.chainInfo }
func (d DaemonInfo) GetConfigFolder() string     { return d.configFolder }
func (d DaemonInfo) GetMoniker() string          { return d.moniker }
func (d DaemonInfo) GetValidatorKeyName() string { return d.validatorKeyName }
func (d DaemonInfo) GetVersion() string          { return d.version }
func (d DaemonInfo) GetPorts() Ports             { return d.ports }
