package commoncmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/playground/cosmosdaemon"
)

const LogLevelFlag = "log_level"

// GetStartOptionsFromCmd parses the command flags to derive the
// start options for the node or chain.
func GetStartOptionsFromCmd(cmd *cobra.Command) (cosmosdaemon.StartOptions, error) {
	ll, err := cmd.Flags().GetString(LogLevelFlag)
	if err != nil {
		return cosmosdaemon.StartOptions{}, fmt.Errorf("could not get %s: %w", LogLevelFlag, err)
	}

	return cosmosdaemon.StartOptions{
		LogLevel: ll,
	}, nil
}
