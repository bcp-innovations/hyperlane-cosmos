package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdSetToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-token [token-id]",
		Short: "Update the Warp token",
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			var ism *util.HexAddress = nil
			if ismId != "" {
				parsed, err := util.DecodeHexAddress(ismId)
				if err != nil {
					return err
				}
				ism = &parsed
			}

			msg := types.MsgSetToken{
				Owner:             clientCtx.GetFromAddress().String(),
				TokenId:           tokenId,
				NewOwner:          newOwner,
				IsmId:             ism,
				RenounceOwnership: renounceOwnership,
			}

			_, err = sdk.AccAddressFromBech32(msg.Owner)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Owner))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringVar(&newOwner, "new-owner", "", "set updated owner")
	cmd.Flags().StringVar(&ismId, "ism-id", "", "set updated ism")
	cmd.Flags().BoolVar(&renounceOwnership, "renounce-ownership", false, "renounce ownership")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
