package keeper_test

import (
	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("hook_igp_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	var mailboxId util.HexAddress
	var hookId util.HexAddress

	const igpDenom = "uatom"
	const actualDenom = "uosmo"

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		mailboxId, err = createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())

		hookId, err = createDummyIgp(s, creator.Address, igpDenom)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: hookId,
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 2,
				GasOracle: &types.GasOracle{
					TokenExchangeRate: math.NewInt(1e10),
					GasPrice:          math.NewInt(1e9),
				},
				GasOverhead: math.NewInt(100000),
			},
		})
		Expect(err).To(BeNil())
	})

	It("IGP PayForGas should fail when any required payment exceeds maxFee", func() {
		recipient, err := util.DecodeHexAddress("0x00000000000000000000000000000000000000000000000000000000deadbeef")
		Expect(err).To(BeNil())

		sender, err := util.DecodeHexAddress("0x0000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496")
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       0,
			Origin:      1,
			Sender:      sender,
			Destination: 2,
			Recipient:   recipient,
			Body:        []byte("test"),
		}

		metadata := util.StandardHookMetadata{
			Address:            creator.AccAddress,
			GasLimit:           math.NewInt(100000),
			CustomHookMetadata: nil,
		}

		// igpDenom != actualDenom
		maxFee := sdk.NewCoins(sdk.NewCoin(actualDenom, math.NewInt(50)))

		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, hookId, metadata, message, maxFee)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("required payment exceeds max hyperlane fee"))
	})
})

func createDummyIgp(s *i.KeeperTestSuite, creator string, denom string) (util.HexAddress, error) {
	res, err := s.RunTx(&types.MsgCreateIgp{
		Owner: creator,
		Denom: denom,
	})
	if err != nil {
		return [32]byte{}, err
	}

	var response types.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	if err != nil {
		return [32]byte{}, err
	}

	return response.Id, nil
}
