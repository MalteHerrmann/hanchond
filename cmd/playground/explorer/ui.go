package explorer

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/explorer"
	"github.com/hanchon/hanchond/playground/explorer/explorerui"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// ui represents the query command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Args:  cobra.ExactArgs(0),
	Short: "Start the node explorer",
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		startingHeight, err := cmd.Flags().GetInt("starting-height")
		if err != nil {
			utils.ExitError(fmt.Errorf("starting height not set"))
		}

		// TODO: move the newFromDB to cosmos daemon
		e := evmos.NewEvmosFromDB(queries, nodeID)
		// TODO: support mainnet and testnet endpoints
		ex := explorer.NewLocalExplorerClient(e.Ports.P8545, e.Ports.P1317, e.HomeDir)

		p := explorerui.CreateExplorerTUI(startingHeight, ex)
		if _, err := p.Run(); err != nil {
			utils.ExitError(fmt.Errorf("error: %w", err))
		}

		utils.ExitSuccess()
	},
}

func init() {
	ExplorerCmd.AddCommand(uiCmd)
	uiCmd.Flags().Int("starting-height", 1, "Starting height to index the chain.")
}
