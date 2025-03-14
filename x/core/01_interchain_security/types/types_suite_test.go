package types_test

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestISMTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/hyperlane/%s Types Test Suite", types.SubModuleName))
}

func hyperlaneMessageFromHexString(hexString string) util.HyperlaneMessage {
	bytes, err := hex.DecodeString(hexString)
	Expect(err).To(BeNil())

	message, err := util.ParseHyperlaneMessage(bytes)
	Expect(err).To(BeNil())

	return message
}

func bytesFromHexString(hexString string) []byte {
	bytes, err := hex.DecodeString(hexString)
	Expect(err).To(BeNil())
	return bytes
}

func signDigest(digest []byte, privateKeyHex string) []byte {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	Expect(err).To(BeNil())

	signature, err := crypto.Sign(digest, privateKey)
	Expect(err).To(BeNil())
	Expect(len(signature)).To(Equal(65))

	signature[64] += 27

	return signature
}
