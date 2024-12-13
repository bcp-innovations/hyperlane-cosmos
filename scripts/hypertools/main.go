package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hypertools",
		Short: "Debug tools for Hyperlane",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "decode-message [hex-encoded-message]",
		Short: "Decodes a Hyperlane message into human readable format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Decode(args[0])
		}})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "warp-transfer [sender-contract] [recipient-contract] [recipient-user] [amount]",
		Short: "Decodes a Hyperlane message into human readable format",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			amount, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			return GenerateWarpTransfer(args[0], args[1], args[2], amount)
		}})

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("failed to execute root command: %w", err))
	}
}
