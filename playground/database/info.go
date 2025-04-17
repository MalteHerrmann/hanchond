package database

import (
	"github.com/hanchon/hanchond/playground/types"
)

// MustParseChainInfo panics if the row's chain info cannot be parsed
// from its string representation that's stored in the database.
func (r GetAllChainNodesRow) MustParseChainInfo() types.ChainInfo {
	return types.MustParseChainInfo(r.ChainInfo)
}

// MustParseChainInfo panics if the chains's info cannot be parsed
// from its string representation that's stored in the database.
func (c Chain) MustParseChainInfo() types.ChainInfo {
	return types.MustParseChainInfo(c.ChainInfo)
}
