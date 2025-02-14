package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

var (
	gasLimit string
	igpId    string
	ismId    string
	maxFee   string
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
		CmdSetIsm(),
		CmdUnrollRemoteRouter(),
	)

	return txCmd
}
