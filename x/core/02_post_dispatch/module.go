package post_dispatch

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/spf13/cobra"
)

import "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/client/cli"

// Name returns the Submodule name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command.
func GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterMsgServer ...
func RegisterMsgServer(server grpc.Server, msgServer types.MsgServer) {
	types.RegisterMsgServer(server, msgServer)
}

// RegisterQueryService registers the gRPC query service.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
