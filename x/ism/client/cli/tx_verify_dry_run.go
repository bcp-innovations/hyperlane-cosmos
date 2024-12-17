package cli

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/ism/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdVerifyDryRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-dry-run [ism-id] [message] [metadata]",
		Short: "Dry-run an ISM verification",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgVerifyDryRun{
				IsmId:    args[0],
				Message:  args[1],
				Metadata: args[2],
				Creator:  clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
