<!--

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

-->

# CHANGELOG

An '!' indicates a state machine breaking change.

## [v1.0.0-beta0](https://github.com/bcp-innovations/hyperlane-cosmos/releases/tag/v1.0.0-beta0) - 2025-03-27

**Initial Release of the Hyperlane Cosmos SDK module**

### Features
- Full mailbox send and receive functionality
- Warp Routes (Token Bridging)
  - Collateral Tokens
  - Synthetic Tokens
  - TokenRouter
- Interchain Security Modules:
  - Merkle-Root-Multisig-ISM
  - MessageId-Multisig-ISM
- Post Dispatch Hooks:
  - Merkle Tree Hook (for Multisig ISMs)
  - InterchainGasPaymaster
