package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func NewIgpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "igp",
		Short: "Hyperlane Interchain Gas Paymaster commands",
	}

	cmd.AddCommand(
		CmdClaim(),
		CmdCreateIgp(),
		CmdSetIgpOwner(),
		CmdPayForGas(),
		CmdSetDestinationGasConfig(),
	)

	return cmd
}

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

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgClaim{
				Sender: clientCtx.GetFromAddress().String(),
				IgpId:  igpId,
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
		Use:   "create [denom]",
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

func CmdSetIgpOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-owner [igp-id] [new-owner]",
		Short: "Update a Hyperlane Interchain Gas Paymaster",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if renounceOwnership && !yes {
				fmt.Print("Are you sure you want to renounce ownership? This action is irreversible. (yes/no): ")
				var response string

				_, err := fmt.Scanln(&response)
				if err != nil {
					return err
				}

				if strings.ToLower(response) != "yes" {
					return fmt.Errorf("canceled transaction")
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgSetIgpOwner{
				Owner:             clientCtx.GetFromAddress().String(),
				IgpId:             igpId,
				NewOwner:          args[1],
				RenounceOwnership: renounceOwnership,
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().BoolVar(&renounceOwnership, "renounce-ownership", false, "renounce ownership")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "auto-confirm all prompts")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdPayForGas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-for-gas [igp-id] [message-id] [destination-domain] [gas-limit] [amount]",
		Short: "Hyperlane Interchain Gas Payment without using QuoteGasPayment",
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

			gasLimitInt, ok := math.NewIntFromString(args[3])
			if !ok {
				return errors.New("failed to convert `gasLimit` into math.Int")
			}

			amount, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			messageId, err := util.DecodeHexAddress(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgPayForGas{
				Sender:            clientCtx.GetFromAddress().String(),
				IgpId:             igpId,
				MessageId:         messageId,
				DestinationDomain: uint32(destinationDomain),
				GasLimit:          gasLimitInt,
				Amount:            amount,
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

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgSetDestinationGasConfig{
				Owner: clientCtx.GetFromAddress().String(),
				IgpId: igpId,
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
