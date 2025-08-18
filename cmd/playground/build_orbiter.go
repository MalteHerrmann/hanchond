package playground

import (
	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/playground/orbiter"
)

var buildOrbiterCmd = &cobra.Command{
	Use:   "build-orbiter",
	Short: "Build the simapp contained in the noble-assets/orbiter repository.",
	Long:  `It downloads, builds and cleans up temporary files for any Orbiter tag. Using the --path flag, you can build a local version of the repository.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		chainInfo := orbiter.ChainInfo

		// TODO: rename if it's working; shouldn't be custom to EVM
		return RunBuildEVMChainCmd(cmd, chainInfo, version)
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildOrbiterCmd)
	// TODO: move this into global flags? Or at least unify with other occurrences
	buildOrbiterCmd.Flags().StringP("path", "p", "", "Path to your local copy of the repository")
}
