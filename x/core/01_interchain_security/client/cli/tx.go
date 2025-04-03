package cli

import (
	"strconv"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
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
		CmdCreateRoutingIsm(),
		CmdSetRoutingIsmDomain(),
		CmdRemoveRoutingIsmDomain(),
	)

	return txCmd
}

func CmdAnnounceValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announce [address] [storage-location] [signature] [mailbox-id]",
		Short: "Announce a Hyperlane validator",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[3])
			if err != nil {
				return err
			}

			msg := types.MsgAnnounceValidator{
				Validator:       args[0],
				StorageLocation: args[1],
				// Expected to be Hex encoded
				Signature: args[2],
				MailboxId: mailboxId,
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
		Use:   "create-message-id-multisig [validators] [threshold]",
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
		Use:   "create-merkle-root-multisig [validators] [threshold]",
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
		Use:   "create-noop",
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

func CmdCreateRoutingIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-routing",
		Short: "Create a Hyperlane Routing ISM",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateRoutingIsm{
				Creator: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetRoutingIsmDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-routing-ism-domain [routing-ism-id] [domain] [ism-id]",
		Short: "Sets the ISM for a given domain in the routing ISM",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			routingIsmId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			domain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			ismId, err := util.DecodeHexAddress(args[2])
			if err != nil {
				return err
			}

			msg := types.MsgSetRoutingIsmDomain{
				IsmId: routingIsmId,
				Route: types.Route{
					Ism:    ismId,
					Domain: uint32(domain),
				},
				Owner: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRemoveRoutingIsmDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-routing-ism-domain [routing-ism-id] [domain]",
		Short: "Removes the ISM for a given domain in the routing ISM",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			routingIsmId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			domain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			msg := types.MsgRemoveRoutingIsmDomain{
				IsmId:  routingIsmId,
				Domain: uint32(domain),
				Owner:  clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
