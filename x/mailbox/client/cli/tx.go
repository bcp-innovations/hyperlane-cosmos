package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

var (
	gasLimit    uint64
	igpId       string
	igpRequired bool
	maxFee      uint64
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
		CmdAnnounceValidator(),
		CmdCreateMailbox(),
		CmdDispatchMessage(),
		CmdProcessMessage(),

		// IGP
		CmdClaim(),
		CmdCreateIgp(),
		CmdPayForGas(),
		CmdSetDestinationGasConfig(),
	)

	return txCmd
}
