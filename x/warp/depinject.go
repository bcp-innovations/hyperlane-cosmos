package warp

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	mailboxkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/keeper"
	mailboxTypes "github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	modulev1 "github.com/bcp-innovations/hyperlane-cosmos/api/warp/module"
)

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService
	AddressCodec address.Codec

	Config *modulev1.Module

	BankKeeper    bankkeeper.Keeper
	MailboxKeeper *mailboxkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	Module appmodule.AppModule
	Keeper keeper.Keeper
	Hooks  mailboxTypes.MailboxHooksWrapper
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance as authority if not provided
	authority := authtypes.NewModuleAddress("gov")
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(in.Cdc, in.AddressCodec, in.StoreService, authority.String(), in.BankKeeper, in.MailboxKeeper)
	m := NewAppModule(in.Cdc, k)
	return ModuleOutputs{Module: m, Keeper: k, Hooks: mailboxTypes.MailboxHooksWrapper{MailboxHooks: k}}
}
