package types

import (
	"fmt"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/ethereum/go-ethereum/crypto"
)

type MultisigISM interface {
	GetValidators() []string
	GetThreshold() uint32
}

func VerifyMultisig(validators []string, threshold uint32, signatures [][]byte, digest [32]byte) (bool, error) {
	if threshold == 0 {
		return false, fmt.Errorf("invalid ism. no threshold present")
	}

	validSignatures := uint32(0)

	// Early return if we can't possibly meet the threshold
	if len(signatures) < int(threshold) {
		return false, nil
	}

	// Get validator public keys
	validatorAddresses := make(map[string]bool, len(validators))
	for _, address := range validators {
		validatorAddresses[strings.ToLower(address)] = true
	}

	for i := 0; i < len(signatures) && validSignatures < threshold; i++ {
		recoveredPubkey, err := util.RecoverEthSignature(digest[:], signatures[i])
		if err != nil {
			continue // Skip invalid signatures
		}

		addressBytes := crypto.PubkeyToAddress(*recoveredPubkey)
		address := util.EncodeEthHex(addressBytes[:])
		if validatorAddresses[address] {
			validSignatures++
		}
	}

	if validSignatures >= threshold {
		return true, nil
	}
	return false, nil
}

func ValidateNewMultisig(m MultisigISM) error {
	if m.GetThreshold() == 0 {
		return fmt.Errorf("threshold must be greater than zero")
	}

	if len(m.GetValidators()) < int(m.GetThreshold()) {
		return fmt.Errorf("validator addresses less than threshold")
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
