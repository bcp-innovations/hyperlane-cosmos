package types

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

type MultisigISM interface {
	GetValidators() []string
	GetThreshold() uint32
}

// VerifyMultisig checks if a message digest is signed by a sufficient number of validators.
// It recovers public keys from signatures and ensures the threshold is met before returning success.
func VerifyMultisig(validators []string, threshold uint32, signatures [][]byte, digest [32]byte) (bool, error) {
	// Check if the number of provided signatures meets the threshold requirement
	if len(signatures) < int(threshold) {
		return false, fmt.Errorf("threshold can not be reached")
	}

	validatorCount := len(validators)
	validatorIndex := 0

	// It is assumed that the signatures are ordered the same way as the validators.
	for i := 0; i < int(threshold); i++ {
		recoveredPubkey, err := util.RecoverEthSignature(digest[:], signatures[i])
		if err != nil {
			return false, fmt.Errorf("failed to recover validator signature: %w", err)
		}

		signerBytes := crypto.PubkeyToAddress(*recoveredPubkey)
		signer := util.EncodeEthHex(signerBytes[:])

		// Loop through remaining validators to find a match for the recovered signer
		for validatorIndex < validatorCount && signer != strings.ToLower(validators[validatorIndex]) {
			// If no match, increment the validator index and continue searching
			validatorIndex++
		}

		// If the validator list was iterated without finding a match, the signature is invalid
		if validatorIndex >= validatorCount {
			return false, nil
		}

		// Move to the next validator for the next signature
		validatorIndex++
	}
	return true, nil
}

// ValidateNewMultisig ensures the Multisig ISM configuration is valid.
func ValidateNewMultisig(m MultisigISM) error {
	if m.GetThreshold() == 0 {
		return fmt.Errorf("threshold must be greater than zero")
	}

	validators := m.GetValidators()
	if len(validators) < int(m.GetThreshold()) {
		return fmt.Errorf("validator addresses less than threshold")
	}

	// Ensure that validators are sorted in ascending order.
	if !slices.IsSorted(validators) {
		return fmt.Errorf("validator addresses are not sorted correctly in ascending order")
	}

	count := map[string]int{}
	for _, validatorAddress := range validators {
		bytes, err := util.DecodeEthHex(validatorAddress)
		if err != nil {
			return fmt.Errorf("invalid validator address: %s", validatorAddress)
		}

		// Ensure that the address is an eth address with 20 bytes.
		if len(bytes) != 20 {
			return fmt.Errorf("invalid validator address: must be 20 bytes")
		}

		// Check for duplications.
		count[validatorAddress]++
		if count[validatorAddress] > 1 {
			return fmt.Errorf("duplicate validator address: %v", validatorAddress)
		}
	}

	return nil
}
