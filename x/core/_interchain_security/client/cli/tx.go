package cli

import (
	"strconv"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "ism",
		Short:                      "Hyperlane Interchain Security Module commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		CmdAnnounceValidator(),
		CmdCreateMessageIdMultisigIsm(),
		CmdCreateMerkleRootMultiSigIsm(),
		CmdCreateNoopIsm(),
	)

	return txCmd
}

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

func CmdCreateMessageIdMultisigIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-message-id-multisig-ism [validators] [threshold]",
		Short: "Create a Hyperlane MessageId Multisig ISM",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			validators := strings.Split(args[0], ",")
			threshold, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateMessageIdMultisigIsm{
				Creator:    clientCtx.GetFromAddress().String(),
				Validators: validators,
				Threshold:  uint32(threshold),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateMerkleRootMultiSigIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-merkle-root-multisig-ism [validators] [threshold]",
		Short: "Create a Hyperlane MerkleRoot Multisig ISM",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			validators := strings.Split(args[0], ",")
			threshold, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateMerkleRootMultisigIsm{
				Creator:    clientCtx.GetFromAddress().String(),
				Validators: validators,
				Threshold:  uint32(threshold),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateNoopIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-noop-ism",
		Short: "Create a Hyperlane Noop ISM",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateNoopIsm{
				Creator: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
