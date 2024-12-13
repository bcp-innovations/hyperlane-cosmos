package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdCreateSyntheticToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-synthetic-token [origin-mailbox] [receiver-domain] [receiver-contract] [ism-id]",
		Short: "Create a Hyperlane Synthetic Token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			domain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			msg := types.MsgCreateSyntheticToken{
				Creator:          clientCtx.GetFromAddress().String(),
				OriginMailbox:    args[0],
				ReceiverDomain:   uint32(domain),
				ReceiverContract: args[2],
				IsmId:            args[3],
			}

			_, err = sdk.AccAddressFromBech32(msg.Creator)
			if err != nil {
				panic(fmt.Errorf("invalid creator address (%s)", msg.Creator))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
