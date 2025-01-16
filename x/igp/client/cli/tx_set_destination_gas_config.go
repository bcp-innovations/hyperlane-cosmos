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

func CmdSetDestinationGasConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-destination-gas-config [igp-id] [remote-domain] [token-exchange-rate] [gas-price] [gas-overhead]",
		Short: "Set Destination Gas Config for Interchain Gas Paymaster",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			remoteDomain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			tokenExchangeRate, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			gasPrice, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return err
			}

			gasOverhead, err := strconv.ParseUint(args[4], 10, 32)
			if err != nil {
				return err
			}

			msg := types.MsgSetDestinationGasConfig{
				Owner: clientCtx.GetFromAddress().String(),
				IgpId: args[0],
				DestinationGasConfig: &types.DestinationGasConfig{
					RemoteDomain: uint32(remoteDomain),
					GasOracle: &types.GasOracle{
						TokenExchangeRate: tokenExchangeRate,
						GasPrice:          gasPrice,
					},
					GasOverhead: gasOverhead,
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
