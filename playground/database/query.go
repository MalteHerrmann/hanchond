package database

import (
	"context"
	"fmt"
)

// GetAllNodesForChainID returns all the nodes for a given chain ID
// and extends the SQL generated query to check if there are zero nodes
// found.
func (q *Queries) GetAllNodesForChainID(ctx context.Context, chainID int64) ([]GetAllChainNodesRow, error) {
	nodes, err := q.GetAllChainNodes(ctx, chainID)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for chain %d", chainID)
	}

	return nodes, nil
}

// Get retrieves the port information for the given standard port id (e.g. 8545 for JSON-RPC).
func (p Port) Get(port uint16) int64 {
	switch port {
	case 1317:
		return p.P1317
	case 6060:
		return p.P6060
	case 6065:
		return p.P6065
	case 8080:
		return p.P8080
	case 8545:
		return p.P8545
	case 8546:
		return p.P8546
	case 9090:
		return p.P9090
	case 9091:
		return p.P9091
	case 26656:
		return p.P26656
	case 26657:
		return p.P26657
	case 26658:
		return p.P26658
	case 26660:
		return p.P26660
	default:
		panic(fmt.Sprintf("returning port %d not supported", port))
	}
}
