package playground

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/solidity"
	"github.com/spf13/cobra"
)

// buildSolcCmd represents the buildSolc command
var buildSolcCmd = &cobra.Command{
	Use:   "build-solc",
	Short: "Build an specific version of Solc",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = filesmanager.SetHomeFolderFromCobraFlags(cmd)
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			utils.ExitError(fmt.Errorf("could not read the version: %w", err))
		}

		isDarwin, err := cmd.Flags().GetBool("is-darwin")
		if err != nil {
			utils.ExitError(fmt.Errorf("could not read the isDarwin: %w", err))
		}

		// Create build folder if needed
		if err := filesmanager.CreateBuildsDir(); err != nil {
			utils.ExitError(fmt.Errorf("could not create build folder: %w", err))
		}

		utils.Log("Downloading solidity from github: %s", version)
		if err := solidity.DownloadSolcBinary(isDarwin, version); err != nil {
			utils.ExitError(err)
		}

		utils.Log("Solc %s is now available", version)
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildSolcCmd)
	buildSolcCmd.PersistentFlags().StringP("version", "v", "0.8.0", "Solc version to download")
	buildSolcCmd.PersistentFlags().Bool("is-darwin", true, "Is the system MacOS arm?")
}
