package evmos

import (
	"context"
	"log"
	"strconv"

	"github.com/hanchon/hanchond/playground/database"
)

// TODO: remove in the long run.
func NewEvmosFromDB(queries *database.Queries, nodeID string) *Evmos {
	nID, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	node, err := queries.GetChainNode(context.Background(), nID)
	if err != nil {
		log.Panic(err)
	}

	ports := node.GetPorts()

	e := NewEvmos(
		node.Moniker,
		node.Version,
		node.ConfigFolder,
		node.ChainID_2,
		node.ValidatorKeyName,
		&ports,
	)
	e.SetValidatorWallet(node.ValidatorKey, node.ValidatorWallet)

	return e
}
