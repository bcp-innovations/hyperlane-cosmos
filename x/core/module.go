package core

import (
	"context"
	"encoding/json"
	"fmt"

	ismmodule "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security"
	ismkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/keeper"
	ismtypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

	pdmodule "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch"
	pdkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	pdtypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/client/cli"
	keeper2 "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasGenesis     = AppModule{}
	_ appmodule.AppModule   = AppModule{}
)

// ConsensusVersion defines the current module consensus version.
const ConsensusVersion = 1

type AppModule struct {
	cdc    codec.Codec
	keeper *keeper2.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper *keeper2.Keeper) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

// Name returns the mailbox module's name.
func (AppModule) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the mailbox module's types on the LegacyAmino codec.
// New modules do not need to support Amino.
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// core
	cdc.RegisterConcrete(&types.MsgCreateMailbox{}, "hyperlane/v1/MsgCreateMailbox", nil)
	cdc.RegisterConcrete(&types.MsgSetMailbox{}, "hyperlane/v1/MsgSetMailbox", nil)
	cdc.RegisterConcrete(&types.MsgProcessMessage{}, "hyperlane/v1/MsgProcessMessage", nil)

	// 01_interchain_security
	cdc.RegisterConcrete(&ismtypes.MsgAnnounceValidator{}, "hyperlane/v1/MsgAnnounceValidator", nil)
	cdc.RegisterConcrete(&ismtypes.MsgCreateMessageIdMultisigIsm{}, "hyperlane/v1/MsgCreateMessageIdMultisigIsm", nil)
	cdc.RegisterConcrete(&ismtypes.MsgCreateMerkleRootMultisigIsm{}, "hyperlane/v1/MsgCreateMerkleRootMultisigIsm", nil)
	cdc.RegisterConcrete(&ismtypes.MsgCreateNoopIsm{}, "hyperlane/v1/MsgCreateNoopIsm", nil)
	cdc.RegisterConcrete(&ismtypes.MsgCreateRoutingIsm{}, "hyperlane/v1/MsgCreateRoutingIsm", nil)
	cdc.RegisterConcrete(&ismtypes.MsgSetRoutingIsmDomain{}, "hyperlane/v1/MsgSetRoutingIsmDomain", nil)
	cdc.RegisterConcrete(&ismtypes.MsgRemoveRoutingIsmDomain{}, "hyperlane/v1/MsgRemoveRoutingIsmDomain", nil)
	cdc.RegisterConcrete(&ismtypes.MsgUpdateRoutingIsmOwner{}, "hyperlane/v1/MsgUpdateRoutingIsmOwner", nil)

	// 02_post_dispatch
	cdc.RegisterConcrete(&pdtypes.MsgClaim{}, "hyperlane/v1/MsgClaim", nil)
	cdc.RegisterConcrete(&pdtypes.MsgCreateIgp{}, "hyperlane/v1/MsgCreateIgp", nil)
	cdc.RegisterConcrete(&pdtypes.MsgCreateMerkleTreeHook{}, "hyperlane/v1/MsgCreateMerkleTreeHook", nil)
	cdc.RegisterConcrete(&pdtypes.MsgCreateNoopHook{}, "hyperlane/v1/MsgCreateNoopHook", nil)
	cdc.RegisterConcrete(&pdtypes.MsgPayForGas{}, "hyperlane/v1/MsgPayForGas", nil)
	cdc.RegisterConcrete(&pdtypes.MsgSetDestinationGasConfig{}, "hyperlane/v1/MsgSetDestinationGasConfig", nil)
	cdc.RegisterConcrete(&pdtypes.MsgSetIgpOwner{}, "hyperlane/v1/MsgSetIgpOwner", nil)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the mailbox module.
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}

	if err := ismtypes.RegisterQueryHandlerClient(context.Background(), mux, ismtypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}

	if err := pdtypes.RegisterQueryHandlerClient(context.Background(), mux, pdtypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces registers interfaces and implementations of the mailbox module.
func (AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)

	// Register Submodules
	ismtypes.RegisterInterfaces(registry)
	pdtypes.RegisterInterfaces(registry)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper2.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper2.NewQueryServerImpl(am.keeper))

	ismmodule.RegisterMsgServer(cfg.MsgServer(), ismkeeper.NewMsgServerImpl(&am.keeper.IsmKeeper))
	ismmodule.RegisterQueryService(cfg.QueryServer(), ismkeeper.NewQueryServerImpl(&am.keeper.IsmKeeper))

	pdmodule.RegisterMsgServer(cfg.MsgServer(), pdkeeper.NewMsgServerImpl(&am.keeper.PostDispatchKeeper))
	pdmodule.RegisterQueryService(cfg.QueryServer(), pdkeeper.NewQueryServerImpl(&am.keeper.PostDispatchKeeper))
}

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.NewGenesisState())
}

// ValidateGenesis performs genesis state validation for the core module.
func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return data.Validate()
}

// InitGenesis performs genesis initialization for the core module.
// It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	if err := am.keeper.InitGenesis(ctx, &genesisState); err != nil {
		panic(fmt.Sprintf("failed to initialize %s genesis state: %v", types.ModuleName, err))
	}

	ismkeeper.InitGenesis(ctx, am.keeper.IsmKeeper, genesisState.IsmGenesis)
	pdkeeper.InitGenesis(ctx, am.keeper.PostDispatchKeeper, genesisState.PostDispatchGenesis)
}

// ExportGenesis returns the exported genesis
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to export %s genesis state: %v", types.ModuleName, err))
	}

	gs.IsmGenesis = ismkeeper.ExportGenesis(ctx, am.keeper.IsmKeeper)
	gs.PostDispatchGenesis = pdkeeper.ExportGenesis(ctx, am.keeper.PostDispatchKeeper)

	return cdc.MustMarshalJSON(gs)
}

// GetTxCmd implements AppModuleBasic interface
func (am AppModule) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd implements AppModuleBasic interface
func (am AppModule) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}
