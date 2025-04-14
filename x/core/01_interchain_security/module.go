package interchain_security

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/client/cli"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the root command for the core ISMs
func GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root command for the core ISMs
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterMsgServer registers the core ism handler for transactions
func RegisterMsgServer(server grpc.Server, msgServer types.MsgServer) {
	types.RegisterMsgServer(server, msgServer)
}

// RegisterQueryService registers the gRPC query service for api queries
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
