package playground

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/gaia"
	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// startNodeCmd represents the startNode command
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

		node, err := queries.GetNode(context.Background(), idNumber)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the node: %w", err))
		}

		chain, err := queries.GetChain(context.Background(), node.ChainID)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the chain: %w", err))
		}

		ci := chain.MustParseChainInfo()

		var pID int
		switch ci.GetBinaryName() {
		case evmos.ChainInfo.GetBinaryName():
			d := evmos.NewEvmos(
				node.Moniker,
				node.Version,
				node.ConfigFolder,
				chain.ChainID,
				node.ValidatorKeyName,
			)
			pID, err = d.Start()
		case gaia.ChainInfo.GetBinaryName():
			d := gaia.NewGaia(
				node.Moniker,
				node.ConfigFolder,
				chain.ChainID,
				node.ValidatorKeyName,
			)
			pID, err = d.Start()
		case sagaos.ChainInfo.GetBinaryName():
			d := sagaos.NewSagaOS(
				node.Moniker,
				node.Version,
				node.ConfigFolder,
				chain.ChainID,
				node.ValidatorKeyName,
			)
			pID, err = d.Start()
		default:
			panic("invalid binary name: " + ci.GetBinaryName())
		}
		if err != nil {
			utils.ExitError(fmt.Errorf("could not start the node: %w", err))
		}
		fmt.Println("Node is running with pID:", pID)

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
