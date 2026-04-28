package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func NewAggregationHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregation",
		Short: "Hyperlane Aggregation Hook commands",
	}

	cmd.AddCommand(
		CmdCreateAggregationHook(),
	)

	return cmd
}

func CmdCreateAggregationHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [hook-id...]",
		Short: "Create a new aggregation hook with the given ordered child hook ids",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			hooks := make([]util.HexAddress, len(args))
			for i, arg := range args {
				hookId, err := util.DecodeHexAddress(arg)
				if err != nil {
					return err
				}
				hooks[i] = hookId
			}

			msg := types.MsgCreateAggregationHook{
				Owner: clientCtx.GetFromAddress().String(),
				Hooks: hooks,
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid sender address (%s)", msg.Owner))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
