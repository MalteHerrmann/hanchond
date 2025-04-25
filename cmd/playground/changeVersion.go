package playground

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// changeVersionCmd represents the changeVersion command
var changeVersionCmd = &cobra.Command{
	Use:   "change-version [id] [version]",
	Args:  cobra.ExactArgs(2),
	Short: "Change the binary version of the given node",
	Long:  `It will update the database entry for the node, you need to manually stop and re-start the node for it to take effect on the running chain.`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		isChainID, _ := cmd.Flags().GetBool("is-chain-id")
		id := args[0]
		idNumber, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not parse the ID: %w", err))
		}
		binaryVersion := strings.TrimSpace(args[1])

		if isChainID {
			// Update all of the chain's nodes
			nodes, err := queries.GetAllNodesForChainID(context.Background(), idNumber)
			if err != nil {
				utils.ExitError(fmt.Errorf("could not get chain nodes: %w", err))
			}

			for _, v := range nodes {
				updateNodeVersion(queries, v.ID, binaryVersion)
			}
		} else {
			// Update just the node
			updateNodeVersion(queries, idNumber, binaryVersion)
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(changeVersionCmd)
	changeVersionCmd.Flags().Bool("is-chain-id", false, "If the flag is yes, it will assume that the ID is the chain ID. If it is set as false, the ID will be used just for the node.")
}

func updateNodeVersion(queries *database.Queries, nodeID int64, version string) {
	_, err := queries.GetNode(context.Background(), nodeID)
	if err != nil {
		utils.ExitError(fmt.Errorf("could not get the node: %w", err))
	}

	err = queries.SetNodeVersion(context.Background(), database.SetNodeVersionParams{
		Version: version,
		ID:      nodeID,
	})
	if err != nil {
		utils.ExitError(fmt.Errorf("could not update the binary version: %w", err))
	}

	utils.Log("Node %d updated to version %s", nodeID, version)
}
