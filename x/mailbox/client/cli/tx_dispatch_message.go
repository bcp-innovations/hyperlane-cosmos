package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

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

			// TODO: Remove, use message-body instead
			messageBody := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgDispatchMessage{
				MailboxId:   mailboxId,
				Sender:      clientCtx.GetFromAddress().String(),
				Destination: uint32(destinationDomain),
				Recipient:   recipient,
				Body:        messageBody,
				IgpId:       igpId,
				GasLimit:    gasLimit,
				MaxFee:      maxFee,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringVar(&igpId, "igp-id", "", "custom InterchainGasPaymaster ID; only used when IGP is not required")

	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 50000, "InterchainGasPayment gas limit (default: 50,000)")

	// TODO: Use default value
	cmd.Flags().Uint64Var(&maxFee, "max-hyperlane-fee", 0, "maximum Hyperlane InterchainGasPayment")
	if err := cmd.MarkFlagRequired("max-hyperlane-fee"); err != nil {
		panic(fmt.Errorf("flag 'max-hyperlane-fee' is required: %w", err))
	}

	return cmd
}
