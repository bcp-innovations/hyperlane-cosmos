package main

import (
	"github.com/spf13/cobra"
)

// Flags
var (
	privateKeys       []string
	senderContract    string
	recipientContract string
	recipientUser     string
	amount            uint64
	message           string
	privateKey        string
	storageLocation   string
	mailboxID         string
	localDomain       uint32
)

func init() {
	rootCmd.AddCommand(announceCmd)
	rootCmd.AddCommand(decodeMessageCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(warpTransferCmd)
}

var rootCmd = &cobra.Command{
	Use:   "hypertools",
	Short: "Debug tools for Hyperlane",
}

func main() {
	rootCmd.Execute()
}
