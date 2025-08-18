package database

import (
	"github.com/hanchon/hanchond/playground/types"
)

// MustParseChainInfo panics if the row's chain info cannot be parsed
// from its string representation that's stored in the database.
func (r GetAllChainNodesRow) MustParseChainInfo() types.ChainInfo {
	return types.MustParseChainInfo(r.ChainInfo)
}

func (r GetAllChainNodesRow) GetDaemonInfo() types.DaemonInfo {
	return types.NewDaemonInfo(
		r.Moniker,
		r.Version,
		r.ConfigFolder,
		r.ChainID_2,
		r.ValidatorKeyName,
		r.MustParseChainInfo(),
		r.GetPorts(),
	)
}

func (r GetAllChainNodesRow) GetPorts() types.Ports {
	return types.Ports{
		P1317:  int(r.P1317),
		P8080:  int(r.P8080),
		P9090:  int(r.P9090),
		P9091:  int(r.P9091),
		P8545:  int(r.P8545),
		P8546:  int(r.P8546),
		P6065:  int(r.P6065),
		P26658: int(r.P26658),
		P26657: int(r.P26657),
		P6060:  int(r.P6060),
		P26656: int(r.P26656),
		P26660: int(r.P26660),
	}
}

// MustParseChainInfo panics if the chains's info cannot be parsed
// from its string representation that's stored in the database.
func (c Chain) MustParseChainInfo() types.ChainInfo {
	return types.MustParseChainInfo(c.ChainInfo)
}

// MustParseChainInfo panics if the node's info cannot be parsed
// from its string representation which is stored in the database.
func (n GetChainNodeRow) MustParseChainInfo() types.ChainInfo {
	return types.MustParseChainInfo(n.ChainInfo)
}

func (n GetChainNodeRow) GetDaemonInfo() types.DaemonInfo {
	return types.NewDaemonInfo(
		n.Moniker,
		n.Version,
		n.ConfigFolder,
		n.ChainID_2,
		n.ValidatorKeyName,
		n.MustParseChainInfo(),
		n.GetPorts(),
	)
}

func (n GetChainNodeRow) GetPorts() types.Ports {
	return types.Ports{
		P1317:  int(n.P1317),
		P8080:  int(n.P8080),
		P9090:  int(n.P9090),
		P9091:  int(n.P9091),
		P8545:  int(n.P8545),
		P8546:  int(n.P8546),
		P6065:  int(n.P6065),
		P26658: int(n.P26658),
		P26657: int(n.P26657),
		P6060:  int(n.P6060),
		P26656: int(n.P26656),
		P26660: int(n.P26660),
	}
}
