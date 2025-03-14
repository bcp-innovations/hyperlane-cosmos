package cli

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdRemoteTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [token-id] [destination-domain] [recipient] [amount]",
		Short: "Send Hyperlane Token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			destinationDomain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			recipient, err := util.DecodeHexAddress(args[2])
			if err != nil {
				return err
			}

			argAmount, ok := math.NewIntFromString(args[3])
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

			maxFeeCoin, err := sdk.ParseCoinNormalized(maxFee)
			if err != nil {
				return err
			}

			var parsedHookId *util.HexAddress = nil
			if customHookId != "" {
				parsed, err := util.DecodeHexAddress(customHookId)
				if err != nil {
					return err
				}
				parsedHookId = &parsed
			}

			msg := types.MsgRemoteTransfer{
				TokenId:            tokenId,
				DestinationDomain:  uint32(destinationDomain),
				Sender:             clientCtx.GetFromAddress().String(),
				Recipient:          recipient,
				Amount:             argAmount,
				CustomHookId:       parsedHookId,
				GasLimit:           gasLimitInt,
				MaxFee:             maxFeeCoin,
				CustomHookMetadata: customHookMetadata,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringVar(&customHookId, "custom-hook-id", "", "custom DefaultHookId")
	cmd.Flags().StringVar(&customHookMetadata, "custom-hook-metadata", "", "custom hook metadata")

	cmd.Flags().StringVar(&gasLimit, "gas-limit", "0", "Overwrite InterchainGasPayment gas limit")

	cmd.Flags().StringVar(&maxFee, "max-hyperlane-fee", "0", "maximum Hyperlane InterchainGasPayment")

	return cmd
}
