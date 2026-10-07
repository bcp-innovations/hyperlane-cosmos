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
* Create (invalid) - non-existing module
* Create (valid)
* Verify (valid) - threshold of modules provided
* Verify (valid) - each module receives its own metadata slice
* Verify (valid) - nested aggregation and routing ISMs
* Verify (valid) - real MessageIdMultisig sub-ISM
* Verify (invalid) - fewer metadatas than threshold
* Verify (invalid) - more metadatas than threshold
* Verify (invalid) - sub-ISM rejects message
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

	It("Create (invalid) - empty modules", func() {
		_, err := createAggregationIsm(nil, 1)
		Expect(err.Error()).To(Equal("modules cannot be empty: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - zero threshold", func() {
		_, err := createAggregationIsm(registerMockIsms(2), 0)
		Expect(err.Error()).To(Equal("threshold must be greater than zero: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - threshold exceeds modules", func() {
		_, err := createAggregationIsm(registerMockIsms(2), 3)
		Expect(err.Error()).To(Equal("threshold 3 exceeds number of modules 2: invalid aggregation ism configuration"))
	})

	It("Create (invalid) - duplicated modules", func() {
		isms := registerMockIsms(1)
		_, err := createAggregationIsm([]util.HexAddress{isms[0], isms[0]}, 2)
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicated module %s: invalid aggregation ism configuration", isms[0].String())))
	})

	It("Create (invalid) - non-existing module", func() {
		isms := registerMockIsms(1)
		unknown := isms[0]
		unknown[31] = 99

		_, err := createAggregationIsm([]util.HexAddress{isms[0], unknown}, 1)
		Expect(err.Error()).To(Equal(fmt.Sprintf("ISM %s not found: unknown ism id", unknown.String())))
	})

	It("Create (valid)", func() {
		modules := registerMockIsms(3)

		ismId := mustCreateAggregationIsm(modules, 2)

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
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)

		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, nil, {0x03}})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		Expect(mockIsm.CallCount()).To(Equal(2))
	})

	It("Verify (valid) - each module receives its own metadata slice", func() {
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)

		metadata := types.FormatAggregationMetadata([][]byte{nil, []byte("second"), []byte("third-metadata")})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

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

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), outer, metadata, util.HyperlaneMessage{Origin: 1})

		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		Expect(mockIsm.CallCount()).To(Equal(2))
		received, _ := mockIsm.ReceivedMetadata(modules[1])
		Expect(received).To(Equal([]byte{0xbb}))
	})

	It("Verify (valid) - real MessageIdMultisig sub-ISM", func() {
		privateKey, err := crypto.GenerateKey()
		Expect(err).To(BeNil())
		validator := crypto.PubkeyToAddress(privateKey.PublicKey)

		res, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{validator.Hex()},
			Threshold:  1,
		})
		Expect(err).To(BeNil())
		var multisigResponse types.MsgCreateMessageIdMultisigIsmResponse
		Expect(proto.Unmarshal(res.MsgResponses[0].Value, &multisigResponse)).To(Succeed())

		mockModules := registerMockIsms(1)
		ismId := mustCreateAggregationIsm([]util.HexAddress{multisigResponse.Id, mockModules[0]}, 2)

		message := util.HyperlaneMessage{Version: 3, Nonce: 7, Origin: 1, Destination: 2, Body: []byte("hello")}
		multisigMetadata := types.MessageIdMultisigMetadata{MerkleIndex: 7}
		digest := multisigMetadata.Digest(&message)
		signature, err := crypto.Sign(digest[:], privateKey)
		Expect(err).To(BeNil())
		signature[64] += 27
		multisigMetadata.Signatures = [][]byte{signature}

		metadata := types.FormatAggregationMetadata([][]byte{multisigMetadata.Bytes(), {0x01}})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, message)
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		// A signature over a different message must fail the whole aggregation.
		otherMessage := message
		otherMessage.Nonce = 8
		result, err = s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, otherMessage)
		Expect(err).To(BeNil())
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - fewer metadatas than threshold", func() {
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)

		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, nil, nil})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err.Error()).To(Equal("expected metadata for 2 modules, got 1: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - more metadatas than threshold", func() {
		// Matches EVM, where providing more metadatas than the threshold reverts.
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 2)

		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, {0x02}, {0x03}})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err.Error()).To(Equal("expected metadata for 2 modules, got 3: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - sub-ISM rejects message", func() {
		modules := registerMockIsms(2)
		rejecting, err := mockIsm.RegisterRejectingIsm(s.Ctx())
		Expect(err).To(BeNil())
		modules = append(modules, rejecting)
		ismId := mustCreateAggregationIsm(modules, 2)

		// The two accepting ISMs alone would satisfy the threshold, but every
		// provided metadata must verify.
		metadata := types.FormatAggregationMetadata([][]byte{{0x01}, nil, {0x03}})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err).To(BeNil())
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - metadata shorter than range header", func() {
		modules := registerMockIsms(3)
		ismId := mustCreateAggregationIsm(modules, 1)

		metadata := types.FormatAggregationMetadata([][]byte{{0x01}})

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err.Error()).To(Equal("invalid metadata length: got 9, expected at least 24 bytes: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - metadata range out of bounds", func() {
		modules := registerMockIsms(1)
		ismId := mustCreateAggregationIsm(modules, 1)

		metadata := []byte{0, 0, 0, 8, 0, 0, 0, 20, 0x01}

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err.Error()).To(Equal("invalid metadata range for module 0: [8:20] with metadata length 9: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - metadata range with start after end", func() {
		modules := registerMockIsms(1)
		ismId := mustCreateAggregationIsm(modules, 1)

		metadata := []byte{0, 0, 0, 9, 0, 0, 0, 8, 0x01}

		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), ismId, metadata, util.HyperlaneMessage{})

		Expect(err.Error()).To(Equal("invalid metadata range for module 0: [9:8] with metadata length 9: invalid aggregation ism metadata"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - gas limit bounds number of sub-ISM calls", func() {
		modules := registerMockIsms(types.MaxAggregationModules)
		ismId := mustCreateAggregationIsm(modules, types.MaxAggregationModules)

		subMetadata := make([][]byte, len(modules))
		for k := range subMetadata {
			subMetadata[k] = []byte{0x01}
		}
		metadata := types.FormatAggregationMetadata(subMetadata)

		// Each Verify call consumes 10_000 gas, so 255 sub-ISMs cannot fit into 1_000_000 gas.
		ctx := s.Ctx().WithGasMeter(storetypes.NewGasMeter(1_000_000))

		Expect(func() {
			_, _ = s.App().HyperlaneKeeper.Verify(ctx, ismId, metadata, util.HyperlaneMessage{})
		}).To(Panic())
	})

	It("Genesis (valid) - export and import of aggregation ISM", func() {
		modules := []util.HexAddress{createNoopIsm(s, creator.Address), createNoopIsm(s, creator.Address)}
		ismId := mustCreateAggregationIsm(modules, 1)

		ismGenesis := keeper.ExportGenesis(s.Ctx(), s.App().HyperlaneKeeper.IsmKeeper)
		coreGenesis, err := s.App().HyperlaneKeeper.ExportGenesis(s.Ctx())
		Expect(err).To(BeNil())

		// import into a fresh chain
		s = i.NewCleanChain()
		Expect(s.App().HyperlaneKeeper.InitGenesis(s.Ctx(), coreGenesis)).To(Succeed())
		keeper.InitGenesis(s.Ctx(), s.App().HyperlaneKeeper.IsmKeeper, ismGenesis)

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
