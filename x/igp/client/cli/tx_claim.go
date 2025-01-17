package cli

import (
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [igp-id]",
		Short: "Claim Hyperlane Interchain Gas Paymaster fees",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgClaim{
				Sender: clientCtx.GetFromAddress().String(),
				IgpId:  args[0],
			}

			_, err = sdk.AccAddressFromBech32(msg.Sender)
			if err != nil {
				panic(fmt.Errorf("invalid sender address (%s)", msg.Sender))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
