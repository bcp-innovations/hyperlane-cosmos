package post_dispatch

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/client/cli"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the root command for the core post dispatch hooks
func GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// RegisterMsgServer registers the post dispatch hook handler for transactions
func RegisterMsgServer(server grpc.Server, msgServer types.MsgServer) {
	types.RegisterMsgServer(server, msgServer)
}

// RegisterQueryService registers the gRPC query service for API queries
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
