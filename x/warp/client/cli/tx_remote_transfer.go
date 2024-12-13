package cli

import (
	"cosmossdk.io/math"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdRemoteTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [token-id] [recipient] [amount]",
		Short: "Send Hyperlane Token",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			tokenId := args[0]
			recipient := args[1]
			argAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return errors.New("invalid amount")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgRemoteTransfer{
				TokenId:   tokenId,
				Sender:    clientCtx.GetFromAddress().String(),
				Recipient: recipient,
				Amount:    argAmount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
