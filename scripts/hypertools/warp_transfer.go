package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	warpTransferCmd.Flags().StringVarP(&senderContract, "sender-contract", "s", "", "Sender contract address")
	if err := warpTransferCmd.MarkFlagRequired("sender-contract"); err != nil {
		panic(fmt.Errorf("failed to mark 'sender-contract' flag as required: %w", err))
	}

	warpTransferCmd.Flags().StringVarP(&recipientContract, "recipient-contract", "r", "", "Recipient contract address")
	if err := warpTransferCmd.MarkFlagRequired("recipient-contract"); err != nil {
		panic(fmt.Errorf("failed to mark 'recipient-contract' flag as required: %w", err))
	}

	warpTransferCmd.Flags().Uint64VarP(&amount, "amount", "a", 0, "Amount of tokens to transfer")
	if err := warpTransferCmd.MarkFlagRequired("amount"); err != nil {
		panic(fmt.Errorf("failed to mark 'amount' flag as required: %w", err))
	}

	warpTransferCmd.Flags().StringVarP(&recipientUser, "recipient-user", "u", "", "Recipient user address")
	if err := warpTransferCmd.MarkFlagRequired("recipient-user"); err != nil {
		panic(fmt.Errorf("failed to mark 'recipient-user' flag as required: %w", err))
	}
}

var warpTransferCmd = &cobra.Command{
	Use:   "warp-transfer",
	Short: "Creates a Warp message for sending tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		return GenerateWarpTransfer(senderContract, recipientContract, recipientUser, amount)
	},
}
