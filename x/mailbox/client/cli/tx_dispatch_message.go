package cli

import (
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
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
