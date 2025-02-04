package main

import (
	"fmt"
	"os"

	"github.com/bcp-innovations/hyperlane-cosmos/tests/simapp"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	simapp.SetAddressPrefixes()

	rootCmd := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "hyperlane", simapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
