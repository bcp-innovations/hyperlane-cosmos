package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

var (
	gasLimit           string
	customHookId       string
	customHookMetadata string
	ismId              string
	maxFee             string
	newOwner           string
	renounceOwnership  bool
	yes                bool
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "hyperlane-transfer",
		Short:                      "hyperlane-transfer transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		CmdCreateCollateralToken(),
		CmdCreateSyntheticToken(),
		CmdEnrollRemoteRouter(),
		CmdRemoteTransfer(),
		CmdSetToken(),
		CmdUnrollRemoteRouter(),
	)

	return txCmd
}
