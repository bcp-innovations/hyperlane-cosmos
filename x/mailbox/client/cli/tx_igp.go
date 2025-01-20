package cli

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [igp-id]",
		Short: "Claim Hyperlane Interchain Gas Paymaster fees",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgClaim{
				Sender: clientCtx.GetFromAddress().String(),
				IgpId:  args[0],
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

func CmdCreateIgp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-igp [denom]",
		Short: "Create a Hyperlane Interchain Gas Paymaster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateIgp{
				Owner: clientCtx.GetFromAddress().String(),
				Denom: args[0],
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

func CmdPayForGas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-for-gas [igp-id] [message-id] [destination-domain] [gas-limit] [max-fee]",
		Short: "Hyperlane Interchain Gas Payment",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			destinationDomain, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			gasLimit, ok := math.NewIntFromString(args[3])
			if !ok {
				return errors.New("failed to convert `gasLimit` into math.Int")
			}

			maxFee, ok := math.NewIntFromString(args[4])
			if !ok {
				return errors.New("failed to convert `maxFee` into math.Int")
			}

			msg := types.MsgPayForGas{
				Sender:            clientCtx.GetFromAddress().String(),
				IgpId:             args[0],
				MessageId:         args[1],
				DestinationDomain: uint32(destinationDomain),
				GasLimit:          &gasLimit,
				MaxFee:            &maxFee,
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

			tokenExchangeRate, ok := math.NewIntFromString(args[2])
			if !ok {
				return errors.New("failed to convert `tokenExchangeRate` into math.Int")
			}

			gasPrice, ok := math.NewIntFromString(args[3])
			if !ok {
				return errors.New("failed to convert `gasPrice` into math.Int")
			}

			gasOverhead, ok := math.NewIntFromString(args[4])
			if !ok {
				return errors.New("failed to convert `gasOverhead` into math.Int")
			}

			msg := types.MsgSetDestinationGasConfig{
				Owner: clientCtx.GetFromAddress().String(),
				IgpId: args[0],
				DestinationGasConfig: &types.DestinationGasConfig{
					RemoteDomain: uint32(remoteDomain),
					GasOracle: &types.GasOracle{
						TokenExchangeRate: &tokenExchangeRate,
						GasPrice:          &gasPrice,
					},
					GasOverhead: &gasOverhead,
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
