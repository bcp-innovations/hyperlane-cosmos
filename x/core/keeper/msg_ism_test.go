package keeper_test

import (
	"fmt"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_ism.go

* Create (valid) Noop ISM
* Create (invalid) Multisig ISM with less pubkeys
* Create (invalid) Multisig ISM with invalid threshold
* Create (invalid) Multisig ISM with invalid validator pubkeys
* Create (valid) Multisig ISM

*/

var _ = Describe("msg_ism.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("Create (valid) Noop ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateNoopIsm{
			Creator: creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateNoopIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		ismId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		ism, _ := s.App().HyperlaneKeeper.Isms.Get(s.Ctx(), ismId.Bytes())
		Expect(ism.Creator).To(Equal(creator.Address))
		Expect(ism.IsmType).To(Equal(types.UNUSED))
		Expect(ism.Ism).To(BeAssignableToTypeOf(&types.Ism_Noop{}))
	})

	It("Create (invalid) Multisig ISM with less pubkeys", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMultisigIsm{
			Creator: creator.Address,
			MultiSig: &types.MultiSigIsm{
				ValidatorPubKeys: []string{},
				Threshold:        2,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal("validator pubkeys less than threshold"))
	})

	It("Create (invalid) Multisig ISM with invalid threshold", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMultisigIsm{
			Creator: creator.Address,
			MultiSig: &types.MultiSigIsm{
				ValidatorPubKeys: []string{},
				Threshold:        0,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero"))
	})

	It("Create (invalid) Multisig ISM with invalid validator pubkeys", func() {
		// Arrange
		validValidatorPubKey := "0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05"
		invalidKeys := []string{
			// one character less
			"0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c0",
			// one character more
			"0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c051",
			// invalid character included (`test`)
			"0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b9793test7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05",
		}

		for _, invalidKey := range invalidKeys {
			// Act
			_, err := s.RunTx(&types.MsgCreateMultisigIsm{
				Creator: creator.Address,
				MultiSig: &types.MultiSigIsm{
					ValidatorPubKeys: []string{
						validValidatorPubKey,
						invalidKey,
					},
					Threshold: 2,
				},
			})

			// Assert
			Expect(err.Error()).To(Equal(fmt.Sprintf("invalid validator pub key: %s", invalidKey)))
		}
	})

	It("Create (valid) Multisig ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateMultisigIsm{
			Creator: creator.Address,
			MultiSig: &types.MultiSigIsm{
				ValidatorPubKeys: []string{
					"0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05",
					"0x0417f57017d748288ccf6341993e47618ce3d4d60614ae09f5149acec191fad3fbca5a8ce4144077948c843ea8e863e3997b6da7a1a6d6c9708f658371430ce06b",
					"0x04ce7edc292d7b747fab2f23584bbafaffde5c8ff17cf689969614441e0527b90015ea9fee96aed6d9c0fc2fbe0bd1883dee223b3200246ff1e21976bdbc9a0fc8",
				},
				Threshold: 2,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateMultisigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		ismId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		ism, _ := s.App().HyperlaneKeeper.Isms.Get(s.Ctx(), ismId.Bytes())
		Expect(ism.Creator).To(Equal(creator.Address))
		Expect(ism.IsmType).To(Equal(types.MESSAGE_ID_MULTISIG))
		Expect(ism.Ism).To(BeAssignableToTypeOf(&types.Ism_MultiSig{}))
	})
})
