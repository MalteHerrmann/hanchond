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

// stopHermesCmd represents the stop-hermes command.
var stopHermesCmd = &cobra.Command{
	Use:   "stop-hermes",
	Short: "Stop the relayer",
	Long:  `It gets the PID from the database and send the kill signal to the process`,
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		relayer, err := queries.GetRelayer(context.Background())
		if err != nil {
			utils.ExitError(fmt.Errorf("relayer is not in the database: %w", err))
		}

		// TODO: check if the process is running checking the PID
		if relayer.IsRunning != 1 {
			utils.ExitError(errors.New("relayer is not running"))
		}

		command := exec.Command( //nolint:gosec
			"kill",
			strconv.FormatInt(relayer.ProcessID, 10),
		)

		out, err := command.CombinedOutput()
		if strings.Contains(strings.ToLower(string(out)), "no such process") {
			utils.Log("relayer is not running, updating the database..")
		} else if err != nil {
			utils.ExitError(fmt.Errorf("could not kill the process: %w", err))
		}

		if err := queries.UpdateRelayer(context.Background(), database.UpdateRelayerParams{
			ProcessID: 0,
			IsRunning: 0,
		}); err != nil {
			utils.ExitError(fmt.Errorf("could not update the relayer database: %w", err))
		}

		utils.Log("relayer is no longer running")
	},
}

func init() {
	PlaygroundCmd.AddCommand(stopHermesCmd)
}
