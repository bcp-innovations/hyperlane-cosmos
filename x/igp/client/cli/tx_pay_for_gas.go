package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
)

func CmdPayForGas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-for-gas [igp-id] [message-id] [destination-domain] [gas-limit]",
		Short: "Hyperlane Interchain Gas Payment",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			destinationDomain, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			gasLimit, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return err
			}

			msg := types.MsgPayForGas{
				Sender:            clientCtx.GetFromAddress().String(),
				IgpId:             args[0],
				MessageId:         args[1],
				DestinationDomain: uint32(destinationDomain),
				GasLimit:          gasLimit,
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
