package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/KYVENetwork/hyperlane-cosmos/x/warp/types"
)

func CmdCreateCollateralToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-collateral-token [origin-mailbox] [origin-denom] [receiver-domain] [receiver-contract] [ism-id]",
		Short: "Create a Hyperlane Collateral Token",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			domain, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			msg := types.MsgCreateCollateralToken{
				Creator:          clientCtx.GetFromAddress().String(),
				OriginMailbox:    args[0],
				OriginDenom:      args[1],
				ReceiverDomain:   uint32(domain),
				ReceiverContract: args[3],
				IsmId:            args[4],
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
