package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/types"
)

func NewDaemonFromDB(queries *database.Queries, nodeID string) (*cosmosdaemon.Daemon, error) {
	// NOTE: we're keeping this here to avoid having to manually parse the node ID everywhere
	id, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node ID: %w", err)
	}

	data, err := cosmosdaemon.GetNodeFromDB(queries, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get node information from db: %w", err)
	}

	println("chain info: ", data.Chain.ChainInfo)

	// we should get the chain info from the DB based on the chain ID and then instantiate the correct NewEvmos, NewCosmos, New...
	var daemon *cosmosdaemon.Daemon

	var chainInfo types.ChainInfo
	err = json.Unmarshal([]byte(data.Chain.ChainInfo), &chainInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chain info: %w", err)
	}
	panic("here without error")

	daemon.RestorePortsFromDB(data.Ports)
	daemon.SetValidatorWallet(data.Node.ValidatorKey, data.Node.ValidatorWallet)

	return daemon, nil
}
