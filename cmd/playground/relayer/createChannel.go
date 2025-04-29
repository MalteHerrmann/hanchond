package relayer

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/hermes"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// represents the createChannelCmd command
var createChannelCmd = &cobra.Command{
	Use:   "create-channel [chain_id] [chain_id]",
	Args:  cobra.ExactArgs(2),
	Short: "Create an IBC channel between two chains. The chains must be previously registered",
	Run: func(cmd *cobra.Command, args []string) {
		_ = sql.InitDBFromCmd(cmd)

		h := hermes.NewHermes()
		utils.Log("Relayer initialized")

		chain1 := args[0]
		chain2 := args[1]

		utils.Log("Calling create channel")
		err := h.CreateChannel(chain1, chain2)
		if err != nil {
			utils.ExitError(fmt.Errorf("error creating channel: %w", err))
		}

		utils.Log("Channel created")
	},
}

func init() {
	RelayerCmd.AddCommand(createChannelCmd)
}
