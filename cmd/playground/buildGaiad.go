package playground

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/gaia"
	"github.com/spf13/cobra"
)

// buildGaiadCmd represents the buildGaiad command
var buildGaiadCmd = &cobra.Command{
	Use:   "build-gaiad",
	Short: "Get the Gaiad binary from the github releases",
	Long:  `It downloads the already built gaiad binary from github, it accepts a version flag to specify any tag. It defaults to: v1.9.0.`,
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

		utils.Log("Downloading gaiad from github: %s", version)
		if err = gaia.GetGaiadBinary(isDarwin, version); err != nil {
			utils.ExitError(fmt.Errorf("could not get gaiad from github: %w", err))
		}

		utils.Log("Gaiad is now available")
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildGaiadCmd)
	buildGaiadCmd.PersistentFlags().StringP("version", "v", "v19.1.0", "Gaiad version to download")
	buildGaiadCmd.PersistentFlags().Bool("is-darwin", true, "Is the system MacOS arm?")
}
