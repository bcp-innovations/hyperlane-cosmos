package cli

import (
	"fmt"
	"strconv"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func NewPausableHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pausable",
		Short: "Hyperlane pausable hook commands",
	}

	cmd.AddCommand(
		CmdCreatePausableHook(),
		CmdSetPausableHookPaused(),
	)

	return cmd
}

func CmdCreatePausableHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [mailbox-id]",
		Short: "Create a new pausable hook for a mailbox",
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

			msg := types.MsgCreatePausableHook{
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

func CmdSetPausableHookPaused() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-paused [hook-id] [true|false]",
		Short: "Set whether a pausable hook rejects dispatches",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			hookId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			paused, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse paused value: %w", err)
			}

			msg := types.MsgSetPausableHookPaused{
				Owner:  clientCtx.GetFromAddress().String(),
				HookId: hookId,
				Paused: paused,
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
