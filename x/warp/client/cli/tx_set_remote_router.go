package cli

import (
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

// TODO: remove
func CmdSetRemoteRouter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-remote-router [token-id] [receiver-domain] [receiver-contract] [gas]",
		Short: "Enroll remote router for a certain token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			receiverDomain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			gas, ok := math.NewIntFromString(args[3])
			if !ok {
				return errors.New("failed to convert `gas` into math.Int")
			}

			msg := types.MsgEnrollRemoteRouter{
				Owner:   clientCtx.GetFromAddress().String(),
				TokenId: args[0],
				RemoteRouter: &types.RemoteRouter{
					ReceiverDomain:   uint32(receiverDomain),
					ReceiverContract: args[2],
					Gas:              gas,
				},
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
