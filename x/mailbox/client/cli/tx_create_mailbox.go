package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

func CmdCreateMailbox() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-mailbox",
		Short: "Create a Hyperlane Mailbox",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateMailbox{
				Creator: clientCtx.GetFromAddress().String(),
			}

			_, err = sdk.AccAddressFromBech32(msg.Creator)
			if err != nil {
				panic(fmt.Errorf("invalid creator address (%s)", msg.Creator))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
