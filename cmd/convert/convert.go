package convert

import (
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/spf13/cobra"
)

// ConvertCmd represents the converter command.
var ConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"c"},
	Short:   "converter utils",
	Long:    `Convert wallets, coins and numbers`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			utils.ExitSuccess()
		}
	},
}
