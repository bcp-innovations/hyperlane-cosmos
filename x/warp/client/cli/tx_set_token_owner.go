package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdSetTokenOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-token-owner [token-id] [new-owner]",
		Short: "Update the Interchain Security Module for a certain token - CAUTION: NEW OWNER IS NOT VERIFIED",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgSetTokenOwner{
				Owner:    clientCtx.GetFromAddress().String(),
				TokenId:  args[0],
				NewOwner: args[1],
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			// TODO: Verify newOwner's validity?
			//_, err = sdk.AccAddressFromBech32(msg.NewOwner)
			//if err != nil {
			//	panic(fmt.Errorf("invalid new owner address (%s)", msg.NewOwner))
			//}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
