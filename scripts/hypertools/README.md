# Hypertools

A collection of debug tools for interacting with the Hyperlane protocol.

## Building

```shell
make build 
```

## Usage

### Message decoding

```shell
./build/hypertools decode-message --message <hex-encoded-message>
```

**Example**

```shell
./build/hypertools decode-message --message 0x0100000003000000011900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24000000001900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e240000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
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
./build/hypertools warp-transfer --sender-contract <sender-contract> --recipient-contract <recipient-contract> --recipient-user <recipient-user> --amount <amount>
```

**sender-contract**: is the contract id from the EVM chain. It must match the receiver-contract field in `http://localhost:1317/hyperlane/warp/v1/tokens`

**recipient-contract**: is the TokenId queried from `http://localhost:1317/hyperlane/warp/v1/tokens` 

**recipient-user**: is the cosmos address of the receiver

**amount**: base amount of the tokens sent

**Example**

```shell
./build/hypertools warp-transfer --sender-contract 0x1900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24 --recipient-contract 0x1900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24 --recipient-user kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd --amount 999
```
```text
0x0100000003000000011900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24000000001900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e240000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
```

### Sign

Signs a message for usage with the MultiSig ISM. Use the `--private-keys` flag to
specify a custom key-set. Otherwise, the default keys will be used.

```shell
./build/hypertools sign --message <message>
```

**Example**
```shell
./build/hypertools sign --message 0x0100000003000000011900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e24000000001900d177c9802c4c78a03c50d08a25344e053bf782913fe4fe45bb4021289e240000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000003e7
```
```text
0x55c70da3bdbe81055a29d85ab37cb6ff9412fe724f156c3beb5b13bec61b9db0352a5b5323ee2cc651e0da0d87cec1773a43f809743e851ea8f990965aada50b0052510886cd24b7f3938e9ecd83322dd6b40a4dbf1474d583f8d46aafd876a5c77605740dc5ae0f1ddacbfcf43b47e43c05b5a5505c33112639a8da0e9038f40701a13bc383dc7a7d94b9d7998c42c2fcdd5da4a4ef1063176b9d371775278502bd322abfec0ba5d9454c549912837a877d4b785fdc2de83eddd80ac33b9643681b00
```


### Announce

Signs a validator announcement digest.

```shell
./build/hypertools announce --private-key <private-key> --storage-location <storage-location> --mailbox-id <mailbox-id> --local-domain <local-domain>
```

**Example**
```shell
./build/hypertools announce --private-key fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a18 --storage-location aws://key.pub2 --mailbox-id 0xe81bf6f262305f49f318d68f33b04866f092ffdb2ecf9c98469b4a8b829f65e4 --local-domain 100
```

```text
0xf25ac9de0d0f65be5699b3549b4874a76b8b1a7cfad8a63134599303b196e474764ee5468679ef08dffeb764bf02ea463f2d93da965f415ddb91fbcb342948a301
```