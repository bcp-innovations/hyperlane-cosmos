# hyperlane-cosmos

This project is an implementation of **Hyperlane for the Cosmos SDK**, designed for a seamless interchain communication following the Hyperlane spec. This allows chains built with the Cosmos SDK to communicate with other blockchains using Hyperlane without relying on CosmWasm.

> [!WARNING]  
> This project is currently under development and not intended to be used in production.

## [x/core](./x/core)
`core` is intended to implement the fundamental functionalities of the Hyperlane 
protocol to dispatch and process messages, which can then be used by 
applications like `warp`. It includes mailboxes and registers hooks as well as
Interchain Security Modules (ISMs) that are implemented in the submodules.

## [x/warp](./x/warp)
`warp` extends the core functionality by enabling token creation and cross-chain 
transfers between chains already connected via Hyperlane. These tokens leverage 
modular security through specific ISMs.

_Both modules can be imported into an CosmosSDK-based chain using [dependency injection](https://docs.cosmos.network/main/build/building-modules/depinject)._

## Building from source

To run all build tools, docker is required. 

```
make all
```

To run the test suite:

```
make test
```

More information can be found in the [Contributing](https://github.com/bcp-innovations/hyperlane-cosmos/blob/main/CONTRIBUTING.md).

## Integrating the module

The integration of the hyperlane module is very similar to the integration
of other Cosmos-SDK modules. The import works via dependency injection.
In this guide we show the integration of the Hyperlane Core module and 
the Warp Module which is used for token bridging.

app.go:
```go
package app 

import (
	// Import the core keeper and warp keeper
    _ "github.com/bcp-innovations/hyperlane-cosmos/x/core"
    hyperlaneKeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
    _ "github.com/bcp-innovations/hyperlane-cosmos/x/warp"
    warpKeeper "github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
)

type App struct {
	// ... other keepers

	// Hyperlane
	HyperlaneKeeper *hyperlaneKeeper.Keeper
	WarpKeeper      warpKeeper.Keeper
}


func New(/* args */) {
	// CosmosSDK code

	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Supply(
				logger,
				appOpts,
			),
		),
		// other keepers

		&app.HyperlaneKeeper,
		&app.WarpKeeper,
	); err != nil {
		// return
	}
	
	// If the module is added during a chain upgrade, tell the store loader
	storeUpgrades := storetypes.StoreUpgrades{
		Added: []string{
			hyperlanetypes.ModuleName,
			warptypes.ModuleName,
		},
	}
	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades))
}
```

app.yaml:
```yaml
modules:
  - name: runtime
    config:
      # omitted properties
      init_genesis: [
        # ... other modules
        hyperlane, warp # add hyperlane and warp here
      ] 
  - name: auth
    config:
      "@type": cosmos.auth.module.v1.Module
      bech32_prefix: hyp
      module_account_permissions:
        # ... other modules
        - account: hyperlane
        - account: warp
          permissions: [ minter, burner ] # give burn and mint permissions to the warp module
  - name: hyperlane
    config:
      "@type": hyperlane.core.module.v1.Module
  - name: warp
    config:
      "@type": hyperlane.warp.module.v1.Module
      enabled_tokens:
        - 1 # Enable Collateral tokens
        - 2 # Enable Synthetic tokens
```


## Deploying a mailbox with hypd

```shell
export HYPD_FLAGS=--home test --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp
```

Create a default Merkle Root Multi-Sig ISM:
```shell
./hypd tx hyperlane ism create-merkle-root-multisig [publiyKey1,...] [threshold] $HYPD_FLAGS
```
http://localhost:1317/hyperlane/v1/isms

Alternatively, one can create a No-Op ISM for testing which will allow every incoming message.

Create a new mailbox. The `ism-id` must be the id of the previously created ISM.
The domain is used for identifying networks. A list of existing domains can
be found here https://docs.hyperlane.xyz/docs/reference/domains.
```shell
/hypd tx hyperlane mailbox create [ism-id] [domain] $HYPD_FLAGS
```
http://localhost:1317/hyperlane/v1/mailboxes

Create a merkle tree hook, which is needed if the recipient used the 
MerkleRootMultisigISM. 
```shell
./hypd tx hyperlane hooks merkle create [mailbox-id] $HYPD_FLAGS
```
http://localhost:1317/hyperlane/v1/merkle_tree_hooks

Create an IGP (Interchain-Gas-Paymaster). 
```shell
./hypd tx hyperlane hooks igp create [denom] $HYPD_FLAGS
```
http://localhost:1317/hyperlane/v1/igps

Set the gas config of the IGP. The config is fully compatible to the Hyperlane
spec and the fee is calculated as follows:
`fee = (gas + gas_overhead) * gas_price * token_exchange_rate / 1e10` 
The resulting unit is the denom specified during the creation. 
```shell
./hypd tx hyperlane hooks igp set-destination-gas-config [igp-id] [remote-domain] [token-exchange-rate] [gas-price] [gas-overhead] $HYPD_FLAGS
```
http://localhost:1317/hyperlane/v1/igps/[igp-id]/destination_gas_configs


Update Mailbox with the newly created hooks:
```shell
./hypd tx hyperlane mailbox set [mailbox-id] --default-hook [igp-hook-id] --required-hook [merkle-tree-hook-id] $HYPD_FLAGS
```

## Deploying a collateral token with hpyd

Create collateral token
```shell
./hypd tx hyperlane-transfer create-collateral-token [mailbox-id] [denom] $HYPD_FLAGS
```

Set ISM for collateral token
```shell
./hypd tx hyperlane-transfer set-token [token-id] --ism-id [ism-id] $HYPD_FLAGS
```

Enroll remote router
```shell
./hypd tx hyperlane-transfer enroll-remote-router [token-id] [destination-domain] [recipient-contract] [gas-required-on-destination-chain] $HYPD_FLAGS
```