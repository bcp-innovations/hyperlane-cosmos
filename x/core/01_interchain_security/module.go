package interchain_security

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/client/cli"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/cosmos-sdk/codec"
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

// RegisterLegacyAminoCodec registers the mailbox module's types on the LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&types.MsgAnnounceValidator{}, "hyperlane/v1/MsgAnnounceValidator", nil)
	cdc.RegisterConcrete(&types.MsgCreateMessageIdMultisigIsm{}, "hyperlane/v1/MsgCreateMessageIdMultisigIsm", nil)
	cdc.RegisterConcrete(&types.MsgCreateMerkleRootMultisigIsm{}, "hyperlane/v1/MsgCreateMerkleRootMultisigIsm", nil)
	cdc.RegisterConcrete(&types.MsgCreateNoopIsm{}, "hyperlane/v1/MsgCreateNoopIsm", nil)
	cdc.RegisterConcrete(&types.MsgCreateRoutingIsm{}, "hyperlane/v1/MsgCreateRoutingIsm", nil)
	cdc.RegisterConcrete(&types.MsgSetRoutingIsmDomain{}, "hyperlane/v1/MsgSetRoutingIsmDomain", nil)
	cdc.RegisterConcrete(&types.MsgRemoveRoutingIsmDomain{}, "hyperlane/v1/MsgRemoveRoutingIsmDomain", nil)
	cdc.RegisterConcrete(&types.MsgUpdateRoutingIsmOwner{}, "hyperlane/v1/MsgUpdateRoutingIsmOwner", nil)
}
