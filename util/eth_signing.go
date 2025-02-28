package util

import (
	"crypto/ecdsa"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
)

// GetEthSigningHash hashes a message according to EIP-191.
//
// The data is a UTF-8 encoded string and will enveloped as follows:
// `"\x19Ethereum Signed Message:\n" + message.length + message` and hashed
// using keccak256.
func GetEthSigningHash(msg []byte) [32]byte {
	const prefix = "\x19Ethereum Signed Message:\n"
	bytes := []byte(fmt.Sprintf("%v%v", prefix, len(msg)))
	bytes = slices.Concat(
		bytes,
		msg,
	)
	return crypto.Keccak256Hash(bytes)
}

// RecoverEthSignature recovers the public key from a given signature
//
// Signatures follows EIP-155 with a recovery id of 27 or 28
func RecoverEthSignature(hash []byte, sig []byte) (*ecdsa.PublicKey, error) {
	// Sub 27 of the recovery id according to this - https://eips.ethereum.org/EIPS/eip-155
	// TODO fix
	sig[64] -= 27

	recoveredPubKey, err := crypto.SigToPub(hash[:], sig)
	if err != nil {
		return nil, err
	}

	return recoveredPubKey, nil
}
