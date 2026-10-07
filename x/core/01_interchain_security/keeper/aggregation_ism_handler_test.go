package keeper_test

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - aggregation_ism_handler_test.go

* Create (invalid) - empty modules
* Create (invalid) - zero threshold
* Create (invalid) - threshold exceeds modules
* Create (invalid) - duplicated modules
* Create (invalid) - too many modules
* Create (invalid) - non-existing module
* Create (valid)
* Verify (valid) - threshold of modules provided
* Verify (valid) - each module receives its own metadata slice
* Verify (valid) - nested aggregation and routing ISMs
* Verify (valid) - real MessageIdMultisig sub-ISM
* Verify (invalid) - real MessageIdMultisig sub-ISM with signature for another message
* Verify (invalid) - non-existing aggregation ISM
* Verify (invalid) - ISM id of aggregation type pointing to another ISM
* Verify (invalid) - fewer metadatas than threshold
* Verify (invalid) - more metadatas than threshold
* Verify (invalid) - sub-ISM rejects message
* Verify (invalid) - sub-ISM returns an error
* Verify (invalid) - metadata shorter than range header
* Verify (invalid) - metadata range out of bounds
* Verify (invalid) - metadata range with start after end
* Verify (invalid) - gas limit bounds number of sub-ISM calls
* Genesis (valid) - export and import of aggregation ISM
*/

var _ = Describe("aggregation_ism_handler.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var mockIsm *i.MockIsm

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		mockIsm = i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
	})

	registerMockIsms := func(count int) []util.HexAddress {
		var isms []util.HexAddress
		for k := 0; k < count; k++ {
			ismId, err := mockIsm.RegisterIsm(s.Ctx())
			Expect(err).To(BeNil())
			isms = append(isms, ismId)
		}
		return isms
	}

	createAggregationIsm := func(modules []util.HexAddress, threshold uint32) (util.HexAddress, error) {
		res, err := s.RunTx(&types.MsgCreateAggregationIsm{
			Creator:   creator.Address,
			Modules:   modules,
			Threshold: threshold,
		})
		if err != nil {
			return util.HexAddress{}, err
		}

		var response types.MsgCreateAggregationIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		return response.Id, nil
	}

	mustCreateAggregationIsm := func(modules []util.HexAddress, threshold uint32) util.HexAddress {
		ismId, err := createAggregationIsm(modules, threshold)
		Expect(err).To(BeNil())
		return ismId
	}

	createMessageIdMultisigIsm := func(validator string) util.HexAddress {
		res, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{validator},
			Threshold:  1,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateMessageIdMultisigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		return response.Id
	}

	It("Create (invalid) - empty modules", func() {
		// Act
		_, err := createAggregationIsm(nil, 1)

		// Assert
		Expect(err.Error()).To(Equal("modules cannot be empty: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - zero threshold", func() {
		// Arrange
		modules := registerMockIsms(2)

		// Act
		_, err := createAggregationIsm(modules, 0)

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - threshold exceeds modules", func() {
		// Arrange
		modules := registerMockIsms(2)

		// Act
		_, err := createAggregationIsm(modules, 3)

		// Assert
		Expect(err.Error()).To(Equal("threshold 3 exceeds number of modules 2: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - duplicated modules", func() {
		// Arrange
		isms := registerMockIsms(1)

		// Act
		_, err := createAggregationIsm([]util.HexAddress{isms[0], isms[0]}, 2)

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicated module %s: invalid aggregation ism configuration", isms[0].String())))
	})

	It("Create (invalid) - too many modules", func() {
		// Arrange
		modules := registerMockIsms(types.MaxAggregationModules + 1)

		// Act
		_, err := createAggregationIsm(modules, 1)

		// Assert
		Expect(err.Error()).To(Equal("too many modules: got 256, max 255: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - non-existing module", func() {
		// Arrange
		isms := registerMockIsms(1)
		unknown := isms[0]
		unknown[31] = 99

		// Act
		_, err := createAggregationIsm([]util.HexAddress{isms[0], unknown}, 1)

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ISM %s not found: unknown ism id", unknown.String())))
	})

	It("Create (valid)", func() {
		// Arrange
		modules := registerMockIsms(3)

		// Act
		ismId := mustCreateAggregationIsm(modules, 2)

		// Assert
		var ism types.AggregationISM
		typeUrl := queryISM(&ism, s, ismId.String())
		Expect(typeUrl).To(Equal("/hyperlane.core.interchain_security.v1.AggregationISM"))
		Expect(ism.Id).To(Equal(ismId))
		Expect(ism.Owner).To(Equal(creator.Address))
		Expect(ism.Modules).To(Equal(modules))
		Expect(ism.Threshold).To(Equal(uint32(2)))
		Expect(ism.ModuleType()).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION))
	})

	It("Verify (valid) - threshold of modules provided", func() {
		// Arrange
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, nil, {0x03}})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		Expect(mockIsm.CallCount()).To(Equal(2))
	})

	It("Verify (valid) - each module receives its own metadata slice", func() {
		// Arrange
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)
		metadata := types.FormatAggregationMetadata([][]byte{nil, []byte("second"), []byte("third-metadata")})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		_, called := mockIsm.ReceivedMetadata(modules[0])
		Expect(called).To(BeFalse())
		received, _ := mockIsm.ReceivedMetadata(modules[1])
		Expect(received).To(Equal([]byte("second")))
		received, _ = mockIsm.ReceivedMetadata(modules[2])
		Expect(received).To(Equal([]byte("third-metadata")))
	})

	It("Verify (valid) - nested aggregation and routing ISMs", func() {
		// Arrange
		modules := registerMockIsms(2)
		inner := mustCreateAggregationIsm(modules, 2)

		res, err := s.RunTx(&types.MsgCreateRoutingIsm{
			Creator: creator.Address,
			Routes:  []types.Route{{Domain: 1, Ism: inner}},
		})
		Expect(err).To(BeNil())
		var routingResponse types.MsgCreateRoutingIsmResponse
		Expect(proto.Unmarshal(res.MsgResponses[0].Value, &routingResponse)).To(Succeed())

		noop := createNoopIsm(s, creator.Address)
		outer := mustCreateAggregationIsm([]util.HexAddress{noop, routingResponse.Id}, 2)

		innerMetadata := types.FormatAggregationMetadata([][]byte{{0xaa}, {0xbb}})
		metadata := types.FormatAggregationMetadata([][]byte{{0x00}, innerMetadata})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), outer, metadata, util.HyperlaneMessage{Origin: 1})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		Expect(mockIsm.CallCount()).To(Equal(2))
		received, _ := mockIsm.ReceivedMetadata(modules[1])
		Expect(received).To(Equal([]byte{0xbb}))
	})

	It("Verify (valid) - real MessageIdMultisig sub-ISM", func() {
		// Arrange
		privateKey, err := crypto.GenerateKey()
		Expect(err).To(BeNil())
		multisigIsm := createMessageIdMultisigIsm(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		ismId := mustCreateAggregationIsm([]util.HexAddress{multisigIsm, registerMockIsms(1)[0]}, 2)

		message := util.HyperlaneMessage{Version: 3, Nonce: 7, Origin: 1, Destination: 2, Body: []byte("hello")}
		multisigMetadata := types.MessageIdMultisigMetadata{MerkleIndex: 7}
		digest := multisigMetadata.Digest(&message)
		signature, err := crypto.Sign(digest[:], privateKey)
		Expect(err).To(BeNil())
		signature[64] += 27
		multisigMetadata.Signatures = [][]byte{signature}

		metadata := types.FormatAggregationMetadata([][]byte{multisigMetadata.Bytes(), {0x01}})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, message)

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
	})

	It("Verify (invalid) - real MessageIdMultisig sub-ISM with signature for another message", func() {
		// Arrange
		privateKey, err := crypto.GenerateKey()
		Expect(err).To(BeNil())
		multisigIsm := createMessageIdMultisigIsm(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		ismId := mustCreateAggregationIsm([]util.HexAddress{multisigIsm, registerMockIsms(1)[0]}, 2)

		signedMessage := util.HyperlaneMessage{Version: 3, Nonce: 7, Origin: 1, Destination: 2, Body: []byte("hello")}
		multisigMetadata := types.MessageIdMultisigMetadata{MerkleIndex: 7}
		digest := multisigMetadata.Digest(&signedMessage)
		signature, err := crypto.Sign(digest[:], privateKey)
		Expect(err).To(BeNil())
		signature[64] += 27
		multisigMetadata.Signatures = [][]byte{signature}

		metadata := types.FormatAggregationMetadata([][]byte{multisigMetadata.Bytes(), {0x01}})

		otherMessage := signedMessage
		otherMessage.Nonce = 8

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, otherMessage)

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - non-existing aggregation ISM", func() {
		// Arrange
		ismId := mustCreateAggregationIsm(registerMockIsms(1), 1)
		ismId[31] = 99

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, []byte{}, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("collections: not found: key '99' of type <nil>"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - ISM id of aggregation type pointing to another ISM", func() {
		// Arrange
		noopIsm := createNoopIsm(s, creator.Address)

		// Same internal id as the Noop ISM, but routed to the aggregation handler.
		ismId := noopIsm
		ismId[23] = types.INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, []byte{}, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ISM %s is not an aggregation ISM: invalid ism type", ismId.String())))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - fewer metadatas than threshold", func() {
		// Arrange
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, nil, nil})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("expected metadata for 2 modules, got 1: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - more metadatas than threshold", func() {
		// Arrange
		// Matches EVM, where providing more metadatas than the threshold reverts.
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, {0x02}, {0x03}})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("expected metadata for 2 modules, got 3: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - sub-ISM rejects message", func() {
		// Arrange
		modules := registerMockIsms(2)
		rejecting, err := mockIsm.RegisterRejectingIsm(s.Ctx())
		Expect(err).To(BeNil())
		modules = append(modules, rejecting)
		ismId := mustCreateAggregationIsm(modules, 2)

		// The two accepting ISMs alone would satisfy the threshold, but every
		// provided metadata must verify.
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, nil, {0x03}})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - sub-ISM returns an error", func() {
		// Arrange
		privateKey, err := crypto.GenerateKey()
		Expect(err).To(BeNil())
		multisigIsm := createMessageIdMultisigIsm(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		ismId := mustCreateAggregationIsm([]util.HexAddress{multisigIsm}, 1)

		// MessageIdMultisig metadata must be at least 68 bytes long.
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("sub-ISM %s at index 0: invalid metadata length: got 1, expected at least 68 bytes", multisigIsm.String())))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - metadata shorter than range header", func() {
		// Arrange
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 1)
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}})

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("invalid metadata length: got 9, expected at least 24 bytes: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - metadata range out of bounds", func() {
		// Arrange
		modules := registerMockIsms(1)
		ismId := mustCreateAggregationIsm(modules, 1)
		metadata := []byte{0, 0, 0, 8, 0, 0, 0, 20, 0x01}

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("invalid metadata range for module 0: [8:20] with metadata length 9: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - metadata range with start after end", func() {
		// Arrange
		modules := registerMockIsms(1)
		ismId := mustCreateAggregationIsm(modules, 1)
		metadata := []byte{0, 0, 0, 9, 0, 0, 0, 8, 0x01}

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("invalid metadata range for module 0: [9:8] with metadata length 9: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - gas limit bounds number of sub-ISM calls", func() {
		// Arrange
		modules := registerMockIsms(types.MaxAggregationModules)
		ismId := mustCreateAggregationIsm(modules, types.MaxAggregationModules)

		subMetadata := make([][]byte, len(modules))
		for k := range subMetadata {
			subMetadata[k] = []byte{0x01}
		}
		metadata := types.FormatAggregationMetadata(subMetadata)

		// Each Verify call consumes 10_000 gas, so 255 sub-ISMs cannot fit into 1_000_000 gas.
		ctx := s.Ctx().WithGasMeter(storetypes.NewGasMeter(1_000_000))

		// Act
		act := func() {
			_, _ = s.App().HyperlaneKeeper.Verify(ctx, ismId, metadata, util.HyperlaneMessage{})
		}

		// Assert
		Expect(act).To(Panic())
	})

	It("Genesis (valid) - export and import of aggregation ISM", func() {
		// Arrange
		modules := []util.HexAddress{createNoopIsm(s, creator.Address), createNoopIsm(s, creator.Address)}
		ismId := mustCreateAggregationIsm(modules, 1)

		// Act
		ismGenesis := keeper.ExportGenesis(s.Ctx(), s.App().HyperlaneKeeper.IsmKeeper)
		coreGenesis, err := s.App().HyperlaneKeeper.ExportGenesis(s.Ctx())
		Expect(err).To(BeNil())

		s = i.NewCleanChain()
		Expect(s.App().HyperlaneKeeper.InitGenesis(s.Ctx(), coreGenesis)).To(Succeed())
		keeper.InitGenesis(s.Ctx(), s.App().HyperlaneKeeper.IsmKeeper, ismGenesis)

		// Assert
		var ism types.AggregationISM
		typeUrl := queryISM(&ism, s, ismId.String())
		Expect(typeUrl).To(Equal("/hyperlane.core.interchain_security.v1.AggregationISM"))
		Expect(ism.Modules).To(Equal(modules))
		Expect(ism.Threshold).To(Equal(uint32(1)))

		metadata := types.FormatAggregationMetadata([][]byte{nil, {}})
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
	})
})
