# hyperlane-cosmos

This project is an implementation of Hyperlane for the Cosmos SDK, designed for a seamless interchain communication following the Hyperlane spec. This allows chains built with the Cosmos SDK to communicate with other blockchains using Hyperlane without relying on CosmWasm.

> [!WARNING]  
> This project is currently under development and not intended to be used in production.

## Modules

### [x/mailbox](./x/mailbox)
### [x/ism](./x/ism)
### [x/warp](./x/warp)

## Usage

_All modules can be imported into an Cosmos SDK chain using [dependency injection](https://docs.cosmos.network/main/build/building-modules/depinject)._

#### Create ISM
```
chaind tx ism create-multisig-ism <validator-pubkeys> <threshold>
```
If the transaction was successfully, you can see the created ISM here: _<api-url>/hyperlane/ism/v1/isms_

#### Create Mailbox
```
chaind tx mailbox create-mailbox
```
If the transaction was successfully, you can see the created Mailbox here: _<api-url>/hyperlane/mailbox/v1/mailboxes_

#### Create Warp Collateral Token
```
chaind tx warp create-collateral-token <origin-mailbox> <origin-denom> <receiver-domain> <receiver-contract> <ism-id> 
```
If the transaction was successfully, you can see the created token here: _<api-url>/hyperlane/warp/v1/tokens_

_For local testing, the same mailbox can be used for sending messages._

#### Transfer tokens using Warp
```
chaind tx warp transfer <token-id> <recipient> <amount>
```
After transferring the token, the `Dispatch` event can be obtained through the block results including the `message_body`.

#### Sign message
```
go run scripts/ism_sign.go <message-body-hex>
```
By default, the script includes 3 private keys that can be used for local testing. However, it is highly recommended to use other keys.

The script will output the next command (without flags).

#### Process message
```
chaind tx mailbox process <metadata> <message>
```
After processing the message, the `Process` event can be obtained through the block results.
