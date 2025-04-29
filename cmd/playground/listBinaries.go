package playground

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

// listBinariesCmd represents the list-binaries command
var listBinariesCmd = &cobra.Command{
	Use:   "list-binaries",
	Short: "List of all the binaries that are available",
	Long:  `List all the binaries using all the files stored in the build folder.`,
	Run: func(cmd *cobra.Command, _ []string) {
		_ = filesmanager.SetHomeFolderFromCobraFlags(cmd)
		versions, err := filesmanager.GetAllEvmosdVersions()
		if err != nil {
			utils.ExitError(fmt.Errorf("could not read files in directory: %w", err))
		}
		fmt.Println("Binaries:")
		for _, v := range versions {
			fmt.Printf("\t- %s\n", v)
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(listBinariesCmd)
}
