package playground

import (
	"github.com/hanchon/hanchond/playground/noble"
	"github.com/spf13/cobra"
)

var buildNobleCmd = &cobra.Command{
	Use:   "build-noble",
	Short: "Build the Noble binary.",
	Long:  `It downloads, builds and cleans up temporary files for any Noble tag. Using the --path flag, you can build a local version of the repository.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		chainInfo := noble.ChainInfo

		return RunBuildEVMChainCmd(cmd, chainInfo, version)
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildNobleCmd)
	// TODO: move this into global flags? Or at least unify with other occurrences
	// TODO: refactor to have build subcommand and then per chain another subcommand where the --path stuff is common for all build subcommands
	buildNobleCmd.Flags().StringP("path", "p", "", "Path to your local copy of the repository")
}
