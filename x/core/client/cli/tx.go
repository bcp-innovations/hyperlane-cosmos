package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	ism "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security"
	pdmodule "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch"
)

var (
	// SetMailbox
	defaultIsm        string
	defaultHook       string
	newOwner          string
	renounceOwnership bool
	requiredHook      string
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewMailboxCmd(),
		ism.GetTxCmd(),
		pdmodule.GetTxCmd(),
	)

	return txCmd
}
