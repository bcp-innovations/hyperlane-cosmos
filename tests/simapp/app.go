package simapp

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/encoding/protojson"

	"sigs.k8s.io/yaml"

	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"

	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warpTypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"

	"github.com/cosmos/cosmos-sdk/x/genutil"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	coremodulev1 "github.com/bcp-innovations/hyperlane-cosmos/api/core/module/v1"
	warpmodulev1 "github.com/bcp-innovations/hyperlane-cosmos/api/warp/module/v1"
	hyperlanekeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	warpkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	genutilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	_ "cosmossdk.io/api/cosmos/tx/config/v1"               // import for side-effects
	_ "github.com/bcp-innovations/hyperlane-cosmos/x/core" // import for side-effects
	_ "github.com/bcp-innovations/hyperlane-cosmos/x/warp" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth"                // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"      // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"                // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/distribution"        // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/genutil"             // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/mint"                // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/staking"             // import for side-effects
)

// DefaultNodeHome default home directories for the application daemon
var DefaultNodeHome string

//go:embed app.yaml
var AppConfigYAML []byte

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

func init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultNodeHome = filepath.Join(dirname, ".hyp")
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// Default Cosmos
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper

	// Hyperlane
	HyperlaneKeeper *hyperlanekeeper.Keeper
	WarpKeeper      warpkeeper.Keeper

	// simulation manager
	sm *module.SimulationManager
}

func DefaultHyperlaneModuleConfigs(enabledTokens []int32) []*appv1alpha1.ModuleConfig {
	return []*appv1alpha1.ModuleConfig{
		{
			Name: warpTypes.ModuleName,
			Config: appconfig.WrapAny(&warpmodulev1.Module{
				EnabledTokens: enabledTokens,
			}),
		},
		{
			Name:   coreTypes.ModuleName,
			Config: appconfig.WrapAny(&coremodulev1.Module{}),
		},
	}
}

// AppConfig returns the default app config.
func AppConfig(hyperlaneModuleConfigs []*appv1alpha1.ModuleConfig) depinject.Config {
	// Manually load YAML file, because we want to inject the custom HyperlaneConfig after that.
	j, err := yaml.YAMLToJSON(AppConfigYAML)
	if err != nil {
		return depinject.Error(err)
	}
	cc := &appv1alpha1.Config{}
	err = protojson.Unmarshal(j, cc)
	if err != nil {
		return depinject.Error(err)
	}

	// Add custom hyperlane configs to app.yaml
	cc.Modules = append(cc.Modules, hyperlaneModuleConfigs...)

	return depinject.Configs(
		appconfig.Compose(cc),
		depinject.Supply(
			// needed for genutil commands
			map[string]module.AppModuleBasic{
				genutilTypes.ModuleName: genutil.NewAppModuleBasic(genutilTypes.DefaultMessageValidator),
			},
		),
	)
}

// NewMiniApp returns a reference to an initialized App.
func NewMiniApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*App, error) {
	return NewMiniAppWithCustomConfig(logger, db, traceStore, loadLatest, appOpts, DefaultHyperlaneModuleConfigs([]int32{1, 2}), baseAppOptions...)
}

// NewMiniAppWithCustomConfig returns a reference to an initialized App.
func NewMiniAppWithCustomConfig(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	hyperlaneConfig []*appv1alpha1.ModuleConfig,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*App, error) {
	var (
		app        = &App{}
		appBuilder *runtime.AppBuilder
	)

	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(hyperlaneConfig),
			depinject.Supply(
				logger,
				appOpts,
			),
		),
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.DistrKeeper,
		&app.ConsensusParamsKeeper,

		&app.HyperlaneKeeper,
		&app.WarpKeeper,
	); err != nil {
		return nil, err
	}

	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	/****  Module Options ****/

	// create the simulation manager and define the order of the modules for deterministic simulations
	// NOTE: this is not required apps that don't use the simulator for fuzz testing transactions
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, make(map[string]module.AppModuleSimulation))
	app.sm.RegisterStoreDecoders()

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}

	return app, nil
}

// LegacyAmino returns App's amino codec.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *App) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}
