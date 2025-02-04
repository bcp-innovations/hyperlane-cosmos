## Hyperlane tests

This directory contains all necessary files for building a dummy binary
with pre-initialized state and providing a full test suite for
the integration tests.

### Dummy binary

Run 
```shell
make build-simapp
```

in the root directory. Then go to `./build` and run

```shell
./hypd init-sample-chain --home <chain-home>
./hypd start --home <chain-home>
```

The chain comes configured with three test accounts which are available
in the `<home>/keyring-test`.


### Integration tests

Import the `tests/integration` package. In your code run

```go
integration.NewCleanChain()
```

to obtain a full test-suite with a pre-configured chain.
