package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	ism "github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

var (
	gasLimit    string
	igpId       string
	igpOptional bool
	maxFee      string
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
		NewAnnounceCmd(),
		NewIgpCmd(),
		NewIsmCmd(),
		NewMailboxCmd(),
		ism.GetTxCmd(),
	)

	return txCmd
}
