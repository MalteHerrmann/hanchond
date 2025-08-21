package playground

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/cmd/playground/common"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/sql"
)

// startChainCmd represents the startChainCmd
//
// TODO: refactor all of the available commands to more layers for easier contextualization.
// e.g. use `h p chain start` and `h p node info` instead of `h p start-chain` and `h p get-node`.
var startChainCmd = &cobra.Command{
	Use:   "start-chain [chain_id]",
	Args:  cobra.ExactArgs(1),
	Short: "Start all the validators of the chain",
	Long:  `Start all the required processes to run the chain`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		chainNumber, err := strconv.Atoi(strings.TrimSpace(args[0]))
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid chain id: %w", err))
		}
		nodes, err := queries.GetAllNodesForChainID(context.Background(), int64(chainNumber))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not find the chain: %w", err))
		}

		startOptions, err := commoncmd.GetStartOptionsFromCmd(cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get start options: %w", err))
		}

		for _, v := range nodes {
			di := v.GetDaemonInfo()
			ports := v.GetPorts()

			node, err := commoncmd.GetDaemonForNode(di, &ports)
			if err != nil {
				utils.ExitError(fmt.Errorf("failed to get node daemon: %w", err))
			}

			pID, err := node.Start(startOptions)
			if err != nil {
				utils.ExitError(fmt.Errorf("could not start node: %w", err))
			}

			utils.Log(fmt.Sprintf("node running with pID: %d", pID))
			err = queries.SetProcessID(context.Background(), database.SetProcessIDParams{
				ProcessID: int64(pID),
				IsRunning: 1,
				ID:        v.ID,
			})
			if err != nil {
				utils.ExitError(fmt.Errorf("could not save the process ID to the db: %w", err))
			}
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(startChainCmd)

	startChainCmd.Flags().
		String(commoncmd.LogLevelFlag, "info", "applied log level for the started nodes")
}
