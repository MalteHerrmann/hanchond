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

// stopNodeCmd represents the stopNode command.
var stopNodeCmd = &cobra.Command{
	Use:   "stop-node id",
	Args:  cobra.ExactArgs(1),
	Short: "Stops a running node with the given ID",
	Long:  `Stops the node using the PID stored in the database`,
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

		if node.IsRunning != 1 {
			utils.ExitError(errors.New("the node is not running"))
		}
		command := exec.Command( //nolint:gosec
			"kill",
			strconv.FormatInt(node.ProcessID, 10),
		)

		out, err := command.CombinedOutput()
		if strings.Contains(strings.ToLower(string(out)), "no such process") {
			utils.Log("process is not running, updating the database..")
		} else if err != nil {
			utils.ExitError(fmt.Errorf("could not kill the process: %w", err))
		}

		if err = queries.SetProcessID(context.Background(), database.SetProcessIDParams{
			ProcessID: 0,
			IsRunning: 0,
			ID:        idNumber,
		}); err != nil {
			utils.ExitError(fmt.Errorf("could not update the database: %w", err))
		}

		utils.Log("node is no longer running")
	},
}

func init() {
	PlaygroundCmd.AddCommand(stopNodeCmd)
}
