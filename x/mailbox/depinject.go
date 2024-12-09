package mailbox

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"fmt"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/keeper"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	"golang.org/x/exp/maps"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	modulev1 "github.com/KYVENetwork/hyperlane-cosmos/api/mailbox/module"
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
		appmodule.Invoke(InvokeSetMailboxHooks),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService
	AddressCodec address.Codec

	Config *modulev1.Module
}

type ModuleOutputs struct {
	depinject.Out

	Module appmodule.AppModule
	Keeper *keeper.Keeper
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance as authority if not provided
	authority := authtypes.NewModuleAddress("gov")
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(in.Cdc, in.AddressCodec, in.StoreService, authority.String())
	m := NewAppModule(in.Cdc, &k)

	return ModuleOutputs{Module: m, Keeper: &k}
}

func InvokeSetMailboxHooks(
	keeper *keeper.Keeper,
	mailboxHooks map[string]types.MailboxHooksWrapper,
) error {
	if keeper != nil && mailboxHooks == nil {
		return nil
	}

	modNames := maps.Keys(mailboxHooks)
	order := modNames
	sort.Strings(order)

	if len(order) != len(modNames) {
		return fmt.Errorf("len(hooks_order: %v) != len(hooks modules: %v)", order, modNames)
	}

	if len(modNames) == 0 {
		return nil
	}

	var multiHooks types.MultiMailboxHooks
	for _, modName := range order {
		hook, ok := mailboxHooks[modName]
		if !ok {
			return fmt.Errorf("can't find mailbox hooks for module %s", modName)
		}

		multiHooks = append(multiHooks, hook)
	}

	keeper.SetHooks(multiHooks)
	return nil
}
