package cli

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
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

			gasLimitInt, ok := math.NewIntFromString(gasLimit)
			if !ok {
				return errors.New("failed to convert `gasLimit` into math.Int")
			}

			maxFeeInt, ok := math.NewIntFromString(maxFee)
			if !ok {
				return errors.New("failed to convert `maxFee` into math.Int")
			}

			msg := types.MsgRemoteTransfer{
				TokenId:   tokenId,
				Sender:    clientCtx.GetFromAddress().String(),
				Recipient: recipient,
				Amount:    argAmount,
				IgpId:     igpId,
				GasLimit:  &gasLimitInt,
				MaxFee:    &maxFeeInt,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringVar(&igpId, "igp-id", "", "custom InterchainGasPaymaster ID; only used when IGP is not required")

	cmd.Flags().StringVar(&gasLimit, "gas-limit", "50000", "InterchainGasPayment gas limit (default: 50,000)")

	// TODO: Use default value
	cmd.Flags().StringVar(&maxFee, "max-hyperlane-fee", "0", "maximum Hyperlane InterchainGasPayment")
	if err := cmd.MarkFlagRequired("max-hyperlane-fee"); err != nil {
		panic(fmt.Errorf("flag 'max-hyperlane-fee' is required: %w", err))
	}

	return cmd
}
