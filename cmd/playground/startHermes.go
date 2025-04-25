package playground

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/hermes"
	localsql "github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// startHermesCmd represents the start-hermes command
var startHermesCmd = &cobra.Command{
	Use:   "start-hermes",
	Short: "Starts the relayer",
	Long:  `The command assumes that the relayer was already built and that there is a channel enabled between 2 chains`,
	Run: func(cmd *cobra.Command, _ []string) {
		queries := localsql.InitDBFromCmd(cmd)
		relayer, err := queries.GetRelayer(context.Background())
		if errors.Is(err, sql.ErrNoRows) {
			if err := queries.InitRelayer(context.Background()); err != nil {
				utils.ExitError(fmt.Errorf("could not init the relayer's database: %w", err))
			}
		}

		// TODO: check if the process is running checking the PID
		if relayer.IsRunning == 1 {
			utils.ExitError(fmt.Errorf("the relayer is already running"))
		}

		pid, err := hermes.NewHermes().Start()
		if err != nil {
			utils.ExitError(fmt.Errorf("could not start the relayer: %w", err))
		}
		utils.Log("Hermes running with PID: %d", pid)

		if err := queries.UpdateRelayer(context.Background(), database.UpdateRelayerParams{
			ProcessID: int64(pid),
			IsRunning: 1,
		}); err != nil {
			utils.ExitError(fmt.Errorf("could not update the relayer database: %w", err))
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(startHermesCmd)
}
