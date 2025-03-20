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


## Deploying a mailbox with hypd

```shell
export HYPD_FLAGS=--home test --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp
```

Create a default ISM:
```shell
./hypd tx hyperlane ism create-merkle-root-multisig [publiyKey1,...] [threshold] $HYPD_FLAGS
```
https://api.korellia.kyve.network/hyperlane/v1/isms

Create a new mailbox
```shell
/hypd tx hyperlane mailbox create [ism-id] [domain] $HYPD_FLAGS
```
https://api.korellia.kyve.network/hyperlane/v1/mailboxes

Create a merkle tree hook
```shell
./hypd tx hyperlane hooks merkle create [mailbox-id] $HYPD_FLAGS
```
https://api.korellia.kyve.network/hyperlane/v1/merkle_tree_hooks

Create an IGP
```shell
./hypd tx hyperlane hooks igp create [denom] $HYPD_FLAGS
```
https://api.korellia.kyve.network/hyperlane/v1/igps

TODO: set gas config
https://api.korellia.kyve.network/hyperlane/v1/igps/0x726f757465725f706f73745f6469737061746368000000040000000000000001/destination_gas_configs


Update Mailbox
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
./hypd tx hyperlane-transfer create-collateral-token [mailbox-id] [denom] $HYPD_FLAGS
```

Enroll gas router
```shell
./hypd tx hyperlane-transfer enroll-remote-router [token-id] [destination-domain] [recipient-contract] [gas-required-on-destination-chain] $HYPD_FLAGS
```