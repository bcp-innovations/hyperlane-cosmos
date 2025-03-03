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
