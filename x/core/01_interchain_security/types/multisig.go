package types

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/ethereum/go-ethereum/crypto"
)

type MultisigISM interface {
	GetValidators() []string
	GetThreshold() uint32
}

func VerifyMultisig(validators []string, threshold uint32, signatures [][]byte, digest [32]byte) (bool, error) {
	if len(signatures) < int(threshold) {
		return false, fmt.Errorf("threshold can not be reached")
	}

	validatorCount := len(validators)
	validatorIndex := 0

	for i := 0; i < int(threshold); i++ {
		recoveredPubkey, err := util.RecoverEthSignature(digest[:], signatures[i])
		if err != nil {
			return false, fmt.Errorf("failed to recover validator signature: %w", err)
		}

		signerBytes := crypto.PubkeyToAddress(*recoveredPubkey)
		signer := util.EncodeEthHex(signerBytes[:])

		for validatorIndex < validatorCount && signer != strings.ToLower(validators[validatorIndex]) {
			validatorIndex++
		}

		if validatorIndex >= validatorCount {
			return false, nil
		}
		validatorIndex++
	}
	return true, nil
}

func ValidateNewMultisig(m MultisigISM) error {
	if m.GetThreshold() == 0 {
		return fmt.Errorf("threshold must be greater than zero")
	}

	if len(m.GetValidators()) < int(m.GetThreshold()) {
		return fmt.Errorf("validator addresses less than threshold")
	}

	// Ensure that validators are sorted in ascending order
	if !slices.IsSorted(m.GetValidators()) {
		return fmt.Errorf("validator addresses are not sorted correctly in ascending order")
	}

	for _, validatorAddress := range m.GetValidators() {
		bytes, err := util.DecodeEthHex(validatorAddress)
		if err != nil {
			return fmt.Errorf("invalid validator address: %s", validatorAddress)
		}

		// ensure that the address is an eth address with 20 bytes
		if len(bytes) != 20 {
			return fmt.Errorf("invalid validator address: must be ethereum address (20 bytes)")
		}
	}

	// check for duplications
	count := map[string]int{}
	for _, address := range m.GetValidators() {
		count[address]++
		if count[address] > 1 {
			return fmt.Errorf("duplicate validator address: %v", address)
		}
	}
	return nil
}
