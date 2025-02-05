package keeper_test

import (
	"fmt"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_mailbox.go

* Create (invalid) Mailbox without default ISM and without IGP
* Create (invalid) Mailbox with invalid default ISM and without IGP
* Create (invalid) Mailbox with non-existing default ISM and without IGP
* Create (invalid) Mailbox with valid default ISM (Noop) and invalid IGP
* Create (invalid) Mailbox with valid default ISM (Multisig) and invalid IGP
* Create (invalid) Mailbox with valid default ISM (Noop) and non-existent IGP
* Create (invalid) Mailbox with valid default ISM (Multisig) and non-existent IGP
* Create (valid) Mailbox with NoopISM and required IGP
* Create (valid) Mailbox with MultisigISM and required IGP
* Create (valid) Mailbox with NoopISM and optional IGP
* Create (valid) Mailbox with MultisigISM and optional IGP

*/

var _ = Describe("msg_mailbox.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	// CreateMailbox() tests
	// ISM
	It("Create (invalid) Mailbox without default ISM and without IGP", func() {
		// Arrange
		defaultIsm := ""

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: defaultIsm,
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: invalid hex address length", defaultIsm)))
	})

	It("Create (invalid) Mailbox with invalid default ISM and without IGP", func() {
		// Arrange
		defaultIsm := "0x1234"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: defaultIsm,
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: invalid hex address length", defaultIsm)))
	})

	It("Create (invalid) Mailbox with non-existing default ISM and without IGP", func() {
		// Arrange
		defaultIsm := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: defaultIsm,
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm)))
	})

	// IGP
	It("Create (invalid) Mailbox with valid default ISM (Noop) and invalid IGP", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId := "0x1234"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp id %s is invalid: invalid hex address length", igpId)))
	})

	It("Create (invalid) Mailbox with valid default ISM (Multisig) and invalid IGP", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId := "0x1234"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp id %s is invalid: invalid hex address length", igpId)))
	})

	It("Create (invalid) Mailbox with valid default ISM (Noop) and non-existent IGP", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp with id %s does not exist", igpId)))
	})

	It("Create (invalid) Mailbox with valid default ISM (Multisig) and non-existent IGP", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp with id %s does not exist", igpId)))
	})

	// Mailbox valid cases
	It("Create (valid) Mailbox with NoopISM and required IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: true,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId, ismId, true)
	})

	It("Create (valid) Mailbox with MultisigISM and required IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createMultisigIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: true,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId, ismId, true)
	})

	It("Create (valid) Mailbox with NoopISM and optional IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId, ismId, false)
	})

	It("Create (valid) Mailbox with MultisigISM and optional IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createMultisigIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId,
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId, ismId, false)
	})

	// DispatchMessage() tests
	// @TODO

	// ProcessMessage() tests
	// @TODO
})

// Utils
func createIgp(s *i.KeeperTestSuite, creator string) string {
	res, err := s.RunTx(&types.MsgCreateIgp{
		Owner: creator,
		Denom: "acoin",
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) string {
	res, err := s.RunTx(&types.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) string {
	res, err := s.RunTx(&types.MsgCreateMultisigIsm{
		Creator: creator,
		MultiSig: &types.MultiSigIsm{
			ValidatorPubKeys: []string{
				"0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05",
				"0x0417f57017d748288ccf6341993e47618ce3d4d60614ae09f5149acec191fad3fbca5a8ce4144077948c843ea8e863e3997b6da7a1a6d6c9708f658371430ce06b",
				"0x04ce7edc292d7b747fab2f23584bbafaffde5c8ff17cf689969614441e0527b90015ea9fee96aed6d9c0fc2fbe0bd1883dee223b3200246ff1e21976bdbc9a0fc8",
			},
			Threshold: 2,
		},
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func verifyNewMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, igpId, ismId string, igpRequired bool) {
	var response types.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, _ := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.Bytes())
	Expect(mailbox.Creator).To(Equal(creator))
	Expect(mailbox.Igp.Id).To(Equal(igpId))
	Expect(mailbox.DefaultIsm).To(Equal(ismId))
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

	if igpRequired {
		Expect(mailbox.Igp.Required).To(BeTrue())
	} else {
		Expect(mailbox.Igp.Required).To(BeFalse())
	}

	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Creator).To(Equal(creator))
}
