package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for 02_post_dispatch
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "hooks",
		Short:                      "Hyperlane Hooks commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
	// TODO add query commands
	)

	return queryCmd
}
