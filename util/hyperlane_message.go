package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	VersionOffset     = 0
	NonceOffset       = 1
	OriginOffset      = 5
	SenderOffset      = 9
	DestinationOffset = 41
	RecipientOffset   = 45
	BodyOffset        = 77
)

// HyperlaneMessage implements the Hyperlane message spec defined here
// https://docs.hyperlane.xyz/docs/reference/libraries/message
type HyperlaneMessage struct {
	Version     uint8
	Nonce       uint32
	Origin      uint32
	Sender      HexAddress
	Destination uint32
	Recipient   HexAddress
	Body        []byte
}

// ParseHyperlaneMessage converts a raw byte slice to HyperlaneMessage struct.
func ParseHyperlaneMessage(raw []byte) (HyperlaneMessage, error) {
	message := HyperlaneMessage{}

	if len(raw) < BodyOffset {
		return message, fmt.Errorf("invalid hyperlane message")
	}

	message.Version = raw[VersionOffset]
	message.Nonce = binary.BigEndian.Uint32(raw[NonceOffset:OriginOffset])
	message.Origin = binary.BigEndian.Uint32(raw[OriginOffset:SenderOffset])
	message.Sender = HexAddress(raw[SenderOffset:DestinationOffset])
	message.Destination = binary.BigEndian.Uint32(raw[DestinationOffset:RecipientOffset])
	message.Recipient = HexAddress(raw[RecipientOffset:BodyOffset])
	message.Body = raw[BodyOffset:]

	return message, nil
}

// Id returns the Keccak256 hash of the raw message bytes which
// corresponds the id of a message according to the Hyperlane spec.
func (msg HyperlaneMessage) Id() HexAddress {
	return HexAddress(crypto.Keccak256(msg.Bytes()))
}

func (msg HyperlaneMessage) Bytes() []byte {
	nonce := make([]byte, 4)
	binary.BigEndian.PutUint32(nonce, msg.Nonce)

	origin := make([]byte, 4)
	binary.BigEndian.PutUint32(origin, msg.Origin)

	destination := make([]byte, 4)
	binary.BigEndian.PutUint32(destination, msg.Destination)

	return slices.Concat(
		[]byte{msg.Version},
		nonce,
		origin,
		msg.Sender.Bytes(),
		destination,
		msg.Recipient.Bytes(),
		msg.Body,
	)
}

// String returns the hex-encoded hyperlane bytes prefixed with "0x"
func (msg HyperlaneMessage) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(msg.Bytes()))
}
