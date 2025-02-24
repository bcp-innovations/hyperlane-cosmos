package keeper_test

import (
	"fmt"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_ism.go

* Create (valid) Noop ISM
* Create (invalid) Multisig ISM with less addresses
* Create (invalid) Multisig ISM with invalid threshold
* Create (invalid) Multisig ISM with duplicate validator addresses
* Create (invalid) Multisig ISM with invalid validator addresses
* Create (valid) Multisig ISM

*/

var _ = Describe("msg_server.go", Ordered, func() {
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
		_, err = util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		//ism, _ := s.App().HyperlaneKeeper.Isms.Get(s.Ctx(), ismId.Bytes())
		//Expect(ism.Creator).To(Equal(creator.Address))
		//Expect(ism.IsmType).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED))
		//Expect(ism.Ism).To(BeAssignableToTypeOf(&types.NoopISM{}))
	})

	It("Create (invalid) Multisig ISM with less address", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultiSigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses less than threshold"))
	})

	It("Create (invalid) Multisig ISM with invalid threshold", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultiSigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  0,
		})

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero"))
	})

	It("Create (invalid) Multisig ISM with duplicate validator addresses", func() {
		// Arrange
		invalidAddress := []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
		}

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultiSigIsm{
			Validators: invalidAddress,
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicate validator address: %v", invalidAddress[0])))
	})

	It("Create (invalid) Multisig ISM with invalid validator address", func() {
		// Arrange
		validValidatorAddress := "0xb05b6a0aa112b61a7aa16c19cac27d970692995e"
		invalidAddress := []string{
			// one character less
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995",
			// one character more
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995ef",
			// invalid character included (`t`)
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995t",
		}

		for _, invalidKey := range invalidAddress {
			// Act
			_, err := s.RunTx(&types.MsgCreateMerkleRootMultiSigIsm{
				Creator: creator.Address,
				Validators: []string{
					validValidatorAddress,
					invalidKey,
				},
				Threshold: 2,
			})

			// Assert
			Expect(err.Error()).To(Equal(fmt.Sprintf("invalid validator address: %s", invalidKey)))
		}
	})

	It("Create (valid) Multisig ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateMerkleRootMultiSigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateMerkleRootMultiSigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		//ismId, err := util.DecodeHexAddress(response.Id)
		//Expect(err).To(BeNil())

		//ism, err := keeper.NewQueryServerImpl(&s.App().IsmKeeper).Ism(s.Ctx(), &types.QueryIsmRequest{Id: ismId.String()})
		//Expect(ism.Ism).To(BeAssignableToTypeOf(&types.MerkleRootMultiSigISM{}))
		//Expect(ism).To(Equal(creator.Address))
		//Expect(ism.IsmType).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG))
	})
})
