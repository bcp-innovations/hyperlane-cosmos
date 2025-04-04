package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

func NewMailboxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailbox",
		Short: "Hyperlane Mailbox commands",
	}

	cmd.AddCommand(
		CmdCreateMailbox(),
		CmdProcessMessage(),
		CmdSetMailbox(),
	)

	return cmd
}

func CmdCreateMailbox() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [default-ism-id] [local-domain]",
		Short: "Create a Hyperlane Mailbox",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			defaultIsm, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse default ism: %v", err)
			}

			localDomain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to parse local domain: %v", err)
			}

			msg := types.MsgCreateMailbox{
				Owner:       clientCtx.GetFromAddress().String(),
				DefaultIsm:  defaultIsm,
				LocalDomain: uint32(localDomain),
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

func CmdProcessMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process [mailbox-id] [metadata] [message]",
		Short: "Process a Hyperlane message",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse mailbox id: %v", err)
			}
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

func parseNullableAddress(address string) (*util.HexAddress, error) {
	if address != "" {
		parsed, err := util.DecodeHexAddress(address)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	}
	return nil, nil
}

func CmdSetMailbox() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [mailbox-id]",
		Short: "Update a Hyperlane Mailbox",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			yes, err := cmd.Flags().GetBool("yes")
			if err != nil {
				return err
			}

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
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			defaultIsmId, err := parseNullableAddress(defaultIsm)
			if err != nil {
				return err
			}

			defaultHookId, err := parseNullableAddress(defaultHook)
			if err != nil {
				return err
			}

			requiredHookId, err := parseNullableAddress(requiredHook)
			if err != nil {
				return err
			}

			msg := types.MsgSetMailbox{
				Owner:             clientCtx.GetFromAddress().String(),
				MailboxId:         mailboxId,
				DefaultIsm:        defaultIsmId,
				DefaultHook:       defaultHookId,
				RequiredHook:      requiredHookId,
				NewOwner:          newOwner,
				RenounceOwnership: renounceOwnership,
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringVar(&defaultIsm, "default-ism", "", "set updated defaultIsm")
	cmd.Flags().StringVar(&defaultHook, "default-hook", "", "set updated defaultHook")
	cmd.Flags().StringVar(&requiredHook, "required-hook", "", "set updated requiredHook")
	cmd.Flags().StringVar(&newOwner, "new-owner", "", "set updated owner")
	cmd.Flags().BoolVar(&renounceOwnership, "renounce-ownership", false, "renounce ownership")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
