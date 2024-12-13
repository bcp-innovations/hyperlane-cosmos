package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
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
	)

	return txCmd
}
