package cli

import (
	"github.com/KYVENetwork/hyperlane-cosmos/x/ism/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

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
