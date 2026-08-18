package types

import (
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

type MultisigISM interface {
	GetValidators() []string
	GetThreshold() uint32
}

// VerifyMultisig reports whether enough validators signed the message digest.
func VerifyMultisig(validators []string, threshold uint32, signatures [][]byte, digest [32]byte) (bool, error) {
	// Check if the number of provided signatures meets the threshold requirement
	if len(signatures) < int(threshold) {
		return false, fmt.Errorf("threshold can not be reached")
	}

	// Decode on every verification to protect against previously stored duplicates.
	validatorAddresses, err := decodeValidators(validators)
	if err != nil {
		return false, fmt.Errorf("invalid multisig validator set: %w", err)
	}

	validatorCount := len(validatorAddresses)
	validatorIndex := 0

	// It is assumed that the signatures are ordered the same way as the validators.
	for i := 0; i < int(threshold); i++ {
		recoveredPubkey, err := util.RecoverEthSignature(digest[:], signatures[i])
		if err != nil {
			return false, fmt.Errorf("failed to recover validator signature: %w", err)
		}

		signer := crypto.PubkeyToAddress(*recoveredPubkey)

		// Loop through remaining validators to find a match for the recovered signer
		for validatorIndex < validatorCount && signer != validatorAddresses[validatorIndex] {
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

	validatorAddresses, err := decodeValidators(validators)
	if err != nil {
		return err
	}

	// Compare decoded addresses so hex casing does not affect their order.
	if !slices.IsSortedFunc(validatorAddresses, common.Address.Cmp) {
		return fmt.Errorf("validator addresses are not sorted correctly in ascending order")
	}

	return nil
}

// decodeValidators rejects malformed and duplicate addresses. Duplicates are
// detected on the decoded bytes rather than the hex string, so case variants of
// one address cannot occupy multiple validator slots.
func decodeValidators(validators []string) ([]common.Address, error) {
	decoded := make([]common.Address, 0, len(validators))
	seen := make(map[common.Address]struct{}, len(validators))

	for _, validator := range validators {
		raw, err := util.DecodeEthHex(validator)
		if err != nil {
			return nil, fmt.Errorf("invalid validator address: %s", validator)
		}

		// Ensure that the address is an eth address with 20 bytes.
		if len(raw) != common.AddressLength {
			return nil, fmt.Errorf("invalid validator address: must be %d bytes", common.AddressLength)
		}

		address := common.Address(raw)

		if _, duplicate := seen[address]; duplicate {
			return nil, fmt.Errorf("duplicate validator address: %s", validator)
		}
		seen[address] = struct{}{}

		decoded = append(decoded, address)
	}

	return decoded, nil
}
