package playground

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// stopHermesCmd represents the stop-hermes command
var stopHermesCmd = &cobra.Command{
	Use:   "stop-hermes",
	Short: "Stop the relayer",
	Long:  `It gets the PID from the database and send the kill signal to the process`,
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		relayer, err := queries.GetRelayer(context.Background())
		if err != nil {
			utils.ExitError(fmt.Errorf("the relayer is not in the database: %w", err))
		}

		// TODO: check if the process is running checking the PID
		if relayer.IsRunning != 1 {
			utils.ExitError(fmt.Errorf("the relayer is not running"))
		}

		command := exec.Command( //nolint:gosec
			"kill",
			fmt.Sprintf("%d", relayer.ProcessID),
		)

		out, err := command.CombinedOutput()
		if strings.Contains(strings.ToLower(string(out)), "no such process") {
			fmt.Println("the relayer is not running, updating the database..")
		} else if err != nil {
			utils.ExitError(fmt.Errorf("could not kill the process: %w", err))
		}

		if err := queries.UpdateRelayer(context.Background(), database.UpdateRelayerParams{
			ProcessID: 0,
			IsRunning: 0,
		}); err != nil {
			utils.ExitError(fmt.Errorf("could not update the relayer database: %w", err))
		}

		fmt.Println("Relayer is no longer running")
	},
}

func init() {
	PlaygroundCmd.AddCommand(stopHermesCmd)
}
