package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func NewNoopHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "noop",
		Short: "Hyperlane Noop Hook commands",
	}

	cmd.AddCommand(
		CmdCreateNoopHook(),
	)

	return cmd
}

func CmdCreateNoopHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new noop hook",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateNoopHook{
				Owner: clientCtx.GetFromAddress().String(),
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
