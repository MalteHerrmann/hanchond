package evmos

import (
	"log"
	"strconv"

	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/database"
)

// TODO: this shouldn't be required anymore after fully refactoring to use the chain config
func NewEvmosFromDB(queries *database.Queries, nodeIDFlagValue string) *Evmos {
	// NOTE: we're keeping this here to avoid having to manually parse the node ID everywhere
	nodeID, err := strconv.ParseInt(nodeIDFlagValue, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	data, err := cosmosdaemon.GetNodeFromDB(queries, nodeID)
	if err != nil {
		log.Panic(err)
	}

	println("chain info: ", data.Chain.ChainInfo)
	// we should get the chain info from the DB based on the chain ID and then instantiate the correct NewEvmos, NewCosmos, New...
	e := NewEvmos(data.Node.Moniker, data.Node.Version, data.Node.ConfigFolder, data.Chain.ChainID, data.Node.ValidatorKeyName)
	e.RestorePortsFromDB(data.Ports)
	e.SetValidatorWallet(data.Node.ValidatorKey, data.Node.ValidatorWallet)
	return e
}
