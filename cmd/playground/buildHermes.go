package playground

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

const DefaultHermesVersion = "v1.10.5"

// buildHermesCmd represents the buildHermes command
var buildHermesCmd = &cobra.Command{
	Use:   "build-hermes",
	Short: "Build the Hermes relayer binary",
	Long:  fmt.Sprintf(`It builds the relayer from source, it accepts a version flag to specify any tag. It defaults to: %s.`, DefaultHermesVersion),
	Run: func(cmd *cobra.Command, _ []string) {
		// TODO: download from release page instead of building from source
		_ = filesmanager.SetHomeFolderFromCobraFlags(cmd)
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			utils.ExitError(fmt.Errorf("could not read the version: %w", err))
		}
		// Clone and build
		if err := filesmanager.CreateTempFolder(version); err != nil {
			utils.ExitError(fmt.Errorf("could not create temp folder: %w", err))
		}

		utils.Log("Cloning hermes version: %s", version)
		if err := filesmanager.GitCloneHermesBranch(version); err != nil {
			utils.ExitError(fmt.Errorf("could not clone the hermes version: %w", err))
		}

		utils.Log("Building hermes...")
		if err := filesmanager.BuildHermes(version); err != nil {
			utils.ExitError(fmt.Errorf("error building hermes: %w", err))
		}

		utils.Log("Moving built binary...")
		if err := filesmanager.SaveHermesBuiltVersion(version); err != nil {
			utils.ExitError(fmt.Errorf("could not move the built binary: %w", err))
		}

		utils.Log("Cleaning up...")
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not remove temp folder: %w", err))
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildHermesCmd)
	buildHermesCmd.PersistentFlags().StringP("version", "v", DefaultHermesVersion, "Hermes version to build")
}
