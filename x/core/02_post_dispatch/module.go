package post_dispatch

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/client/cli"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the root command for the core post dispatch hooks
func GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root command for the core post dispatch hooks
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterMsgServer registers the post dispatch hook handler for transactions
func RegisterMsgServer(server grpc.Server, msgServer types.MsgServer) {
	types.RegisterMsgServer(server, msgServer)
}

// RegisterQueryService registers the gRPC query service for API queries
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}

// RegisterLegacyAminoCodec registers the mailbox module's types on the LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&types.MsgClaim{}, "hyperlane/v1/MsgClaim", nil)
	cdc.RegisterConcrete(&types.MsgCreateIgp{}, "hyperlane/v1/MsgCreateIgp", nil)
	cdc.RegisterConcrete(&types.MsgCreateMerkleTreeHook{}, "hyperlane/v1/MsgCreateMerkleTreeHook", nil)
	cdc.RegisterConcrete(&types.MsgCreateNoopHook{}, "hyperlane/v1/MsgCreateNoopHook", nil)
	cdc.RegisterConcrete(&types.MsgPayForGas{}, "hyperlane/v1/MsgPayForGas", nil)
	cdc.RegisterConcrete(&types.MsgSetDestinationGasConfig{}, "hyperlane/v1/MsgSetDestinationGasConfig", nil)
	cdc.RegisterConcrete(&types.MsgSetIgpOwner{}, "hyperlane/v1/MsgSetIgpOwner", nil)
}
