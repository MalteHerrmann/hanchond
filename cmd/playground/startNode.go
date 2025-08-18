package playground

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/cmd/playground/common"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/sql"
)

// startNodeCmd represents the startNode command.
var startNodeCmd = &cobra.Command{
	Use:   "start-node [node_id]",
	Args:  cobra.ExactArgs(1),
	Short: "Starts a node with the given ID",
	Long:  `It will run the node in a subprocess, saving the pid in the database in case it needs to be stopped in the future`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		id := args[0]
		idNumber, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not parse the ID: %w", err))
		}

		node, err := queries.GetChainNode(context.Background(), idNumber)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get chain node: %w", err))
		}

		ports := node.GetPorts()
		d, err := common.GetDaemonForNode(node.GetDaemonInfo(), &ports)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get daemon info: %w", err))
		}

		pID, err := d.Start()
		if err != nil {
			utils.ExitError(fmt.Errorf("could not start the node: %w", err))
		}
		utils.Log("Node is running with pID: %d", pID)

		err = queries.SetProcessID(context.Background(), database.SetProcessIDParams{
			ProcessID: int64(pID),
			IsRunning: 1,
			ID:        node.ID,
		})
		if err != nil {
			utils.ExitError(fmt.Errorf("could not save the process ID to the db: %w", err))
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(startNodeCmd)
}
