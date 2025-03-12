package types_test

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server.go

* Test (valid) multi-sig signature
* Test (invalid) empty metadata
* Test (invalid) empty message
* Test (invalid) multi-sig signature for wrong validator
* Test (invalid) multi-sig signature with wrong threshold
* Test (invalid) duplicated signature

*/

var _ = Describe("msg_server.go", Ordered, func() {
	validMetadata := bytesFromHexString("000000000000000000000000b7f8bc63bbcad18155201308c8f3540b07f84f5e000000007d444379286585e7899dfa4e9ee5687d893dace3cef71f79e882703d52f17dc70000000000000000000000000000000000000000000000000000000000000000ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5b4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d3021ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85e58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a193440eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968ffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f839867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756afcefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0f9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5f8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf8923490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99cc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8beccda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d22733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981fe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0b46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0c65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2f4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd95a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e3774df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652cdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618db8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea32293237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d7358448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a900000000aeae4828232950be5882dc4a7dc3f87c6f7524b09693f622191915dfd47982250aae8da55df59ef2b8fe3c1a08122a2e37babee29f8356e67eee96f40ce2c5031c")

	validMessage := hyperlaneMessageFromHexString("030000000000007a690000000000000000000000004a679253410272dd5232b3ff7cf5dbb88f29531904861f2e726f757465725f617070000000000000000000000000000100000000000000000000000000000000000000009022fae177099ff75c2010db21e05da50bcb109100000000000000000000000000000000000000000000000000000000000f4240")

	It("Test (valid) multi-sig signature", func() {
		// Arrange
		var validators []string
		for _, validator := range PrivateKeys {
			validators = append(validators, validator.address)
		}

		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MerkleRootMultisigMetadata{
			MerkleTreeHook:  [32]byte{},
			MessageIndex:    uint32(0),
			MerkleProof:     [32][32]byte{},
			SignedIndex:     uint32(0),
			SignedMessageId: message.Id(),
		}

		digest := metadata.Digest(&message)

		var signatures [][]byte
		for i := range PrivateKeys {
			sig := signDigest(digest[:], PrivateKeys[i].privateKey)
			signatures = append(signatures, sig)
		}
		metadata.Signatures = signatures

		// Act
		verify, err := merkleRootMultiSig.Verify(sdk.Context{}, metadata.Bytes(), message)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeTrue())
	})

	It("Test (invalid) unsorted validators", func() {
		// Arrange
		var validators []string
		for i := range PrivateKeys {
			validators = append(validators, PrivateKeys[len(PrivateKeys)-(i+1)].address)
		}

		// Act
		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}
		// Assert
		Expect(merkleRootMultiSig.Validate().Error()).To(Equal("validator addresses are not sorted correctly in ascending order"))
	})

	It("Test (invalid) empty metadata", func() {
		// Arrange
		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{"0x0c60e7eCd06429052223C78452F791AAb5C5CAc6"},
			Threshold:  1,
		}

		metadata := bytesFromHexString("")

		// Act
		verify, err := merkleRootMultiSig.Verify(sdk.Context{}, metadata, validMessage)

		// Assert
		Expect(err.Error()).To(Equal("invalid metadata length: got 0, expected at least 1096 bytes"))
		Expect(verify).To(BeFalse())
	})

	It("Test (invalid) empty message", func() {
		// Arrange
		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{"0x0c60e7eCd06429052223C78452F791AAb5C5CAc6"},
			Threshold:  1,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}
		// Act
		verify, err := merkleRootMultiSig.Verify(sdk.Context{}, validMetadata, message)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeFalse())
	})

	It("Test (invalid) multi-sig signature for wrong validator", func() {
		// Arrange
		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{"0x0c60e7eCd06429052223C78452F791AAb5C5CAc7"},
			Threshold:  1,
		}

		// Act
		verify, err := merkleRootMultiSig.Verify(sdk.Context{}, validMetadata, validMessage)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeFalse())
	})

	It("Test (invalid) multi-sig signature with wrong threshold", func() {
		// Arrange
		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{"0x0c60e7eCd06429052223C78452F791AAb5C5CAc6"},
			Threshold:  2,
		}

		// Act
		verify, err := merkleRootMultiSig.Verify(sdk.Context{}, validMetadata, validMessage)

		// Assert
		Expect(err.Error()).To(Equal("threshold can not be reached"))
		Expect(verify).To(BeFalse())
	})

	It("Test (invalid) duplicated signature", func() {
		// Arrange
		var validators []string
		for _, validator := range PrivateKeys {
			validators = append(validators, validator.address)
		}

		merkleRootMultiSig := types.MerkleRootMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MerkleRootMultisigMetadata{
			MerkleTreeHook:  [32]byte{},
			MessageIndex:    uint32(0),
			MerkleProof:     [32][32]byte{},
			SignedIndex:     uint32(0),
			SignedMessageId: message.Id(),
		}

		digest := metadata.Digest(&message)

		var duplicatedSignatures [][]byte
		for range PrivateKeys {
			sig := signDigest(digest[:], PrivateKeys[0].privateKey)
			duplicatedSignatures = append(duplicatedSignatures, sig)
		}
		metadata.Signatures = duplicatedSignatures

		// Act
		verify, err := merkleRootMultiSig.Verify(sdk.Context{}, metadata.Bytes(), message)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeFalse())
	})
})

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
