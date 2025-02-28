package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for accessing the API through the CLI
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "ism",
		Short:                      "ISM client query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
	// TODO(low priority) add query commands
	)

	return queryCmd
}
