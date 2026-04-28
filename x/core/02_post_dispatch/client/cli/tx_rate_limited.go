package cli

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func NewRateLimitedHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate-limit",
		Short: "Hyperlane rate limited hook commands",
	}

	cmd.AddCommand(
		CmdCreateRateLimitedHook(),
		CmdSetRateLimit(),
		CmdRemoveRateLimit(),
	)

	return cmd
}

func CmdCreateRateLimitedHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [mailbox-id]",
		Short: "Create a new rate limited hook for a mailbox",
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

			msg := types.MsgCreateRateLimitedHook{
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

func CmdSetRateLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [hook-id] [token-id] [max-capacity]",
		Short: "Set a token rate limit on a rate limited hook",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			hookId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[1])
			if err != nil {
				return err
			}

			maxCapacity, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("failed to parse max capacity: %s", args[2])
			}

			msg := types.MsgSetRateLimit{
				Owner:       clientCtx.GetFromAddress().String(),
				HookId:      hookId,
				TokenId:     tokenId,
				MaxCapacity: maxCapacity,
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

func CmdRemoveRateLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [hook-id] [token-id]",
		Short: "Remove a token rate limit from a rate limited hook",
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

			tokenId, err := util.DecodeHexAddress(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgRemoveRateLimit{
				Owner:   clientCtx.GetFromAddress().String(),
				HookId:  hookId,
				TokenId: tokenId,
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
