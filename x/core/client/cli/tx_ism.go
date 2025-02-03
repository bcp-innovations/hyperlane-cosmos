package cli

import (
	"strconv"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func NewIsmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ism",
		Short: "Hyperlane Interchain Security Module commands",
	}

	cmd.AddCommand(
		CmdCreateMultiSigIsm(),
		CmdCreateNoopIsm(),
	)

	return cmd
}

func CmdCreateMultiSigIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-multisig-ism [validators] [threshold]",
		Short: "Create a Hyperlane MultiSig ISM",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			validators := strings.Split(args[0], ",")
			threshold, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			multisigIsm := types.MultiSigIsm{
				ValidatorPubKeys: validators,
				Threshold:        uint32(threshold),
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateMultisigIsm{
				Creator:  clientCtx.GetFromAddress().String(),
				MultiSig: &multisigIsm,
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
