package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
)

func CmdCreateIgp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-igp [beneficiary]",
		Short: "Create a Hyperlane Interchain Gas Paymaster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateIgp{
				Owner:       clientCtx.GetFromAddress().String(),
				Beneficiary: args[0],
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			_, err = sdk.AccAddressFromBech32(msg.Beneficiary)
			if err != nil {
				panic(fmt.Errorf("invalid beneficiary address (%s)", msg.Beneficiary))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
