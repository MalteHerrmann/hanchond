package cosmosdaemon

import (
	"context"
	"fmt"

	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/filesmanager"
)

func (d *Daemon) SaveChainToDB(queries *database.Queries) (database.Chain, error) {
	return queries.InsertChain(context.Background(), database.InsertChainParams{
		Name:          fmt.Sprintf("chain-%s", d.ChainID),
		ChainID:       d.ChainID,
		BinaryVersion: filesmanager.GetVersionedBinaryName(d.chainInfo, d.Version),
		Denom:         d.chainInfo.GetDenom(),
		Prefix:        d.chainInfo.GetAccountPrefix(),
	})
}

func (d *Daemon) SaveNodeToDB(chain database.Chain, queries *database.Queries) (int64, error) {
	nodeID, err := queries.InsertNode(context.Background(), database.InsertNodeParams{
		ChainID:          chain.ID,
		ConfigFolder:     d.HomeDir,
		Moniker:          d.Moniker,
		ValidatorKey:     d.ValMnemonic,
		ValidatorKeyName: d.ValKeyName,
		ValidatorWallet:  d.ValWallet,
		KeyType:          string(d.chainInfo.GetKeyAlgo()),
		BinaryVersion:    d.BinaryName,

		ProcessID:   0,
		IsValidator: 1,
		IsArchive:   0,
		IsRunning:   0,
	})
	if err != nil {
		return 0, err
	}

	err = queries.InsertPorts(context.Background(), database.InsertPortsParams{
		NodeID: nodeID,
		P1317:  int64(d.Ports.P1317),
		P8080:  int64(d.Ports.P8080),
		P9090:  int64(d.Ports.P9090),
		P9091:  int64(d.Ports.P9091),
		P8545:  int64(d.Ports.P8545),
		P8546:  int64(d.Ports.P8546),
		P6065:  int64(d.Ports.P6065),
		P26658: int64(d.Ports.P26658),
		P26657: int64(d.Ports.P26657),
		P6060:  int64(d.Ports.P6060),
		P26656: int64(d.Ports.P26656),
		P26660: int64(d.Ports.P26660),
	})
	if err != nil {
		return 0, err
	}

	return nodeID, nil
}
