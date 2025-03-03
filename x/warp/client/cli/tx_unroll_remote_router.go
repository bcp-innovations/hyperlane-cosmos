package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdUnrollRemoteRouter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unroll-remote-router [token-id] [receiver-domain]",
		Short: "Unroll remote router for a certain token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			receiverDomain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			msg := types.MsgUnrollRemoteRouter{
				Owner:          clientCtx.GetFromAddress().String(),
				TokenId:        tokenId,
				ReceiverDomain: uint32(receiverDomain),
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
