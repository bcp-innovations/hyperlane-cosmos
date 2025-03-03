# Contributing

Thank you for considering to contribute to this project. We mostly follow 
Cosmos-SDK principles and design architectures.

## Overview

- The latest state of development is on `main`.
- `main` must always pass `make all`.
- Releases can be found in `/release/*` branch.
- Everything must be covered by tests. We have a very extensive test-suite
  and use triple-A testing (Arrange, Act, Assert).

## Creating a Pull Request

- Check out the latest state from main and always keep the PR in sync with main.
- Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#specification).
- Only one feature per pull request.
- Write an entry for the Changelog.
- Write tests covering 100% of your modified code.
- The command `make all` must pass.

## Coding Guidelines

- Write readable and maintainable code. `Premature Optimization Is the Root of All Evil`.
  Concentrate on clean interfaces first and only optimize for performance if it is needed.
- The keeper directory is structured the following:
    - `msg_server_*`-files are the entry point for message handling. This file
      should be very clean to read and outsource most of the part to the logic files.
      One should immediately understand the flow by just reading the function names
      which are called while handling the message.
    - `query_server_*`-files are the entry point for query handling. This file
      should be very clean to read. Most queries can make usa of the pagination
      provided in `./util`

## Code Structure

- `./util` contains generic functions for handling encodings and global custom types
- `./x/core` contains the CosmosSDK logic for managing a mailbox and coordinates
  the registration of ISMs, PostDispatchHooks and Apps.
  Every module providing additional functionality must register itself in the
  core module.
  - `./x/core/01_interchain_security` provides the basic Hyperlane ISMs ready
    for usage
  - `./x/core/02_post_dispatch` provides basic Hyperlane post-dispatch hooks
    ready for usage.
- `./x/warp` is an external app using the Hyperlane core module. It provides
  functionality for Collateral and Synthetic tokens.

## Legal

You agree that your contribution is licenced under the given LICENSE and all
ownership is handed over to the authors named in 
[LICENSE](https://github.com/bcp-innovations/hyperlane-cosmos/blob/main/LICENSE).
