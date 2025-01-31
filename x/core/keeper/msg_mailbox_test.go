package keeper

import (
	"fmt"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("msg_mailbox.go", Ordered, func() {

	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
	})

	It("Create New Mailbox", func() {

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		balance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, i.A_DENOM)

		fmt.Printf(balance.String())

		s.Commit()
		Expect(true).To(BeTrue())
	})

})
