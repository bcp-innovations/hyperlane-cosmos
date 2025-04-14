package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
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
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%v transactions subcommands", types.ModuleName),
		Aliases:                    []string{"hyperlane-transfer"},
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
