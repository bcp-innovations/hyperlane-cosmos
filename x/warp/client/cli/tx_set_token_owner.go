package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdSetToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-token [token-id]",
		Short: "Update the Warp token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			var ism *util.HexAddress = nil
			if ismId != "" {
				parsed, err := util.DecodeHexAddress(ismId)
				if err != nil {
					return err
				}
				ism = &parsed
			}

			msg := types.MsgSetToken{
				Owner:    clientCtx.GetFromAddress().String(),
				TokenId:  tokenId,
				NewOwner: newOwner,
				IsmId:    ism,
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringVar(&newOwner, "new-owner", "", "set updated owner")
	cmd.Flags().StringVar(&ismId, "ism-id", "", "set updated ism")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
