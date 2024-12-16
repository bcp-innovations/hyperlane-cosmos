# Hypertools

A collection of debug tools for interacting with the Hyperlane protocol.

## Building

```shell
make build 
```

## Usage

### Message decoding

```shell
./build/hypertools decode-message [hex-encoded message]
```

**Example**

```shell
./build/hypertools decode-message 0x0100000003000000011900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24000000001900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e240000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
```
```text
### Message ###
Version:        1
Nonce:          3
Origin:         1
Sender:         0x1900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24
Destination:    0
Recipient:      0x1900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24
Body:           0x0000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
```

### Warp Transfer

Used to generate generic warp transfer messages. The explanation and example
will cover the case when tokens are sent from an EVM chain to a Cosmos chain.

```shell
./build/hypertools warp-transfer [sender-contract] [recipient-contract] [recipient-user] [amount]
```

**sender-contract**: is the contract id from the EVM chain. It must match the receiver-contract field in `http://localhost:1317/hyperlane/warp/v1/tokens`

**recipient-contract**: is the TokenId queried from `http://localhost:1317/hyperlane/warp/v1/tokens` 

**recipient-user**: is the cosmos address of the receiver

**amount**: base amount of the tokens sent

**Example**

```shell
./build/hypertools warp-transfer 0x1900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24 0x1900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24 kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd 999
```
```text
0x0100000003000000011900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24000000001900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e240000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
```

### Sign

Signs a message for usage with the MultiSig ISM. Use the `--private-keys` flag to
specify a custom key-set. Otherwise, the default keys will be used.

```shell
./build/hypertools sign [message]
```

**Example**
```shell
./build/hypertools sign 0x0100000003000000011900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24000000001900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e240000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
```
```text
0x55c70da3bdbe81055a29d85ab37cb6ff9412fe724f156c3beb5b13bec61b9db0352a5b5323ee2cc651e0da0d87cec1773a43f809743e851ea8f990965aada50b0052510886cd24b7f3938e9ecd83322dd6b40a4dbf1474d583f8d46aafd876a5c77605740dc5ae0f1ddacbfcf43b47e43c05b5a5505c33112639a8da0e9038f40701a13bc383dc7a7d94b9d7998c42c2fcdd5da4a4ef1063176b9d371775278502bd322abfec0ba5d9454c549912837a877d4b785fdc2de83eddd80ac33b9643681b00
```
