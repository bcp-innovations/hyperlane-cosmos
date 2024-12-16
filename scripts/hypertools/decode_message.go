package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	decodeMessageCmd.Flags().StringVarP(&message, "message", "m", "", "Hex-encoded message to decode")
	if err := decodeMessageCmd.MarkFlagRequired("message"); err != nil {
		panic(fmt.Errorf("failed to mark 'message' flag as required: %w", err))
	}
}

var decodeMessageCmd = &cobra.Command{
	Use:   "decode-message",
	Short: "Decodes a Hyperlane message into human-readable format",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Decode(message)
	},
}
