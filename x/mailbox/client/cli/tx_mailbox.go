package cli

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

func CmdAnnounceValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announce-validator [address] [storage-location] [signature] [mailbox-id]",
		Short: "Announce a Hyperlane validator",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgAnnounceValidator{
				Validator:       args[0],
				StorageLocation: args[1],
				// Expected to be Hex encoded
				Signature: args[2],
				MailboxId: args[3],
				Creator:   clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateMailbox() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-mailbox [igp-id]",
		Short: "Create a Hyperlane Mailbox",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateMailbox{
				Creator: clientCtx.GetFromAddress().String(),
				Igp: &types.InterchainGasPaymaster{
					Id:       args[0],
					Required: igpRequired,
				},
			}

			_, err = sdk.AccAddressFromBech32(msg.Creator)
			if err != nil {
				panic(fmt.Errorf("invalid creator address (%s)", msg.Creator))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().BoolVar(&igpRequired, "igp-required", true, " determine whether InterchainGasPaymaster is required")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDispatchMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dispatch [mailbox-id] [recipient] [destination-domain] [message-body]",
		Short: "Dispatch a Hyperlane message",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			mailboxId := args[0]
			recipient := args[1]

			destinationDomain, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			messageBody := args[3]

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

			msg := types.MsgDispatchMessage{
				MailboxId:   mailboxId,
				Sender:      clientCtx.GetFromAddress().String(),
				Destination: uint32(destinationDomain),
				Recipient:   recipient,
				Body:        messageBody,
				IgpId:       igpId,
				GasLimit:    &gasLimitInt,
				MaxFee:      &maxFeeInt,
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

func CmdProcessMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process [mailboxId] [metadata] [message]",
		Short: "Process a Hyperlane message",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			mailboxId := args[0]
			metadata := args[1]
			message := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgProcessMessage{
				MailboxId: mailboxId,
				Metadata:  metadata,
				Message:   message,
				Relayer:   clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
