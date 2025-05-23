package playground

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// stopChainCmd represents the stopChain command
var stopChainCmd = &cobra.Command{
	Use:   "stop-chain [chain_id]",
	Args:  cobra.ExactArgs(1),
	Short: "Stops all the running validators for the given Chain ID",
	Long:  `Stops all the nodes using the PID stored in the database`,
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

		for _, v := range nodes {
			if v.IsRunning != 1 {
				utils.Log("The node %d is not running", v.ID)
				continue
			}

			command := exec.Command( //nolint:gosec
				"kill",
				fmt.Sprintf("%d", v.ProcessID),
			)
			out, err := command.CombinedOutput()

			if strings.Contains(strings.ToLower(string(out)), "no such process") {
				utils.Log("Process is not running for node %d, updating the database..", v.ID)
			} else if err != nil {
				utils.ExitError(fmt.Errorf("could not kill the process: %w", err))
			}

			if err = queries.SetProcessID(context.Background(), database.SetProcessIDParams{
				ProcessID: 0,
				IsRunning: 0,
				ID:        v.ID,
			}); err != nil {
				utils.ExitError(fmt.Errorf("could not update the database: %w", err))
			}
		}

		utils.Log("Chain %d is stopped", chainNumber)
	},
}

func init() {
	PlaygroundCmd.AddCommand(stopChainCmd)
}
