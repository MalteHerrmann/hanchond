package convert

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/converter"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/spf13/cobra"
)

// AddrCmd represents the addr command
var AddrCmd = &cobra.Command{
	Use:   "addr",
	Args:  cobra.ExactArgs(1),
	Short: "Convert between bech32 and hex addresses",
	Long:  `Convert between cosmos and ethereum encoded addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		if converter.Has0xPrefix(input) {
			prefix, err := cmd.Flags().GetString("prefix")
			if err != nil {
				utils.ExitError(err)
			}
			value, err := converter.HexToBech32(input, prefix)
			if err != nil {
				utils.ExitError(err)
			}

			fmt.Println(value)
		} else {
			addr, err := converter.Bech32ToHex(input)
			if err != nil {
				utils.ExitError(err)
			}

			fmt.Println(addr)
		}
	},
}

func init() {
	ConvertCmd.AddCommand(AddrCmd)
	AddrCmd.Flags().StringP("prefix", "p", "evmos", "bech32 prefix")
}
