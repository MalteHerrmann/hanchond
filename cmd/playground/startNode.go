package playground

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// startNodeCmd represents the startNode command
var startNodeCmd = &cobra.Command{
	// TODO: this should not require the chain ID but just get the node from the DB should contain all info
	Use:   "start-node [chain_id] [node_id]",
	Args:  cobra.ExactArgs(2),
	Short: "Start a node for a specific chain",
	Long:  `It will start a node for a specific chain. The node will be started with the given chain ID and node ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		chainID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid chain ID"))
		}

		nodeID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid node ID"))
		}

		queries := sql.InitDBFromCmd(cmd)

		chain, err := queries.GetChain(context.Background(), chainID)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the chain info from db: %w", err))
		}

		chainInfo := chain.MustParseChainInfo()
		path := filesmanager.GetNodeHomeFolder(chainID, nodeID)

		daemon := cosmosdaemon.NewDameon(
			chainInfo,
			fmt.Sprintf("moniker-%d-%d", chainID, nodeID),
			// TODO: the version here is wrong, this is passing the whole chainInfo as a string
			chain.ChainInfo,
			path,
			chain.ChainID,
			fmt.Sprintf("validator-key-%d-%d", chainID, nodeID),
		)

		pid, err := daemon.StartNodeAndStoreInfo(queries, nodeID)
		if err != nil {
			utils.ExitError(fmt.Errorf("error starting node: %w", err))
		}

		utils.Log("Node started successfully with PID: %d", pid)
	},
}

func init() {
	PlaygroundCmd.AddCommand(startNodeCmd)
}
