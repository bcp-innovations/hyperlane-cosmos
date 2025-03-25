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

func NewMerkleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merkle",
		Short: "Hyperlane Merkle Tree Hook commands",
	}

	cmd.AddCommand(
		CmdCreateMerkle(),
	)

	return cmd
}

func CmdCreateMerkle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [mailbox-id]",
		Short: "Create a new merkle tree hook with the given mailbox id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgCreateMerkleTreeHook{
				Owner:     clientCtx.GetFromAddress().String(),
				MailboxId: mailboxId,
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
