package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = ""

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the tool version.",
	Args:  cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
