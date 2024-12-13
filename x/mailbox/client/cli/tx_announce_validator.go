package cli

import (
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
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
