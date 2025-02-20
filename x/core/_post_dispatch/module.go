package _interchain_security_module

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/spf13/cobra"
)

import "github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/client/cli"

// Name returns the IBC channel ICS name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for IBC channels.
func GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for IBC channels.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterMsgServer ...
func RegisterMsgServer(server grpc.Server, queryServer types.MsgServer) {
	types.RegisterMsgServer(server, queryServer)
}

//// RegisterQueryService registers the gRPC query service for IBC channels.
//func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
//	types.RegisterQueryServer(server, queryServer)
//}
