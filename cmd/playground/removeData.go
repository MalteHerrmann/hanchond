package playground

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// removeDataCmd represents the removeData command
var removeDataCmd = &cobra.Command{
	Use:   "remove-data",
	Short: "Removes the data folder, deleting the configuration and data for all the networks and relayers",
	Long:  `It is a command useful when restarting the process from scratch, it will delete all the data keeping just the built binaries. NOTE: it will also stop all running services`,
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)

		// Stop all nodes
		fmt.Println("Stopping all the running nodes...")
		stopping := false
		if nodes, err := queries.GetAllNodes(context.Background()); err == nil {
			// Database is initialized
			for _, node := range nodes {
				if node.IsRunning == 1 {
					stopping = true
					command := exec.Command( //nolint:gosec
						"kill",
						fmt.Sprintf("%d", node.ProcessID),
					)
					_, _ = command.CombinedOutput()
				}
			}
		}

		// Stop the relayer
		fmt.Println("Stopping the relayer...")
		if relayer, err := queries.GetRelayer(context.Background()); err == nil {
			// The relayer is runnning
			if relayer.IsRunning == 1 {
				stopping = true
				command := exec.Command( //nolint:gosec
					"kill",
					fmt.Sprintf("%d", relayer.ProcessID),
				)
				_, _ = command.CombinedOutput()
			}
		}

		// If we killed a process, wait 2 secods before deleting the files so the directory is not being used
		if stopping {
			time.Sleep(2 * time.Second)
		}

		// Clean up disk data
		fmt.Println("Cleaning up the data...")
		if err := filesmanager.CleanUpData(); err != nil {
			utils.ExitError(fmt.Errorf("failed to remove the data: %w", err))
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(removeDataCmd)
}
