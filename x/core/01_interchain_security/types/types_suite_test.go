package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestISMTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/hyperlane/%s Types Test Suite", types.SubModuleName))
}

func bytesFromHexString(hexString string) []byte {
	bytes, err := hex.DecodeString(hexString)
	Expect(err).To(BeNil())
	return bytes
}

func hyperlaneMessageFromHexString(hexString string) util.HyperlaneMessage {
	bytes, err := hex.DecodeString(hexString)
	Expect(err).To(BeNil())

	message, err := util.ParseHyperlaneMessage(bytes)
	Expect(err).To(BeNil())

	return message
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

type keypair struct {
	address    string
	privateKey string
}

var PrivateKeys = []keypair{
	{
		address:    "0x06CE2a5ECDc3a0850978664c44327E80E10aF8Ab",
		privateKey: "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a18",
	},
	{
		address:    "0xdf738d27Da985BDdE29E6a34C0a945ff81Aa21DA",
		privateKey: "c87509a1c067bbde78beb793e6fa49e3462d7e7bcfbb4e3f79e926fc27ae42c4",
	},
	{
		address:    "0xf17f52151EbEF6C7334FAD080c5704D77216b732",
		privateKey: "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f",
	},
}
