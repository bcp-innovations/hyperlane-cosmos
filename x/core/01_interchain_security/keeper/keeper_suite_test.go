package keeper_test

import (
	"fmt"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/keeper"
	"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - routing_ism_handler_test.go

* Verify (invalid) non-existing ISM

*/

func TestISMKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/hyperlane/%s Keeper Test Suite", types.SubModuleName))
}

func queryISM(ism proto.Message, s *i.KeeperTestSuite, ismId string) string {
	queryServer := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper)
	rawIsm, err := queryServer.Ism(s.Ctx(), &types.QueryIsmRequest{Id: ismId})
	Expect(err).To(BeNil())
	err = proto.Unmarshal(rawIsm.Ism.Value, ism)
	Expect(err).To(BeNil())
	return rawIsm.Ism.TypeUrl
}

var _ = Describe("msg_server.go", Ordered, func() {
	var s *i.KeeperTestSuite

	BeforeEach(func() {
		s = i.NewCleanChain()
	})

	It("Verify (invalid) non-existing ISM", func() {
		// Arrange
		// Act
		verify, err := s.App().HyperlaneKeeper.IsmKeeper.Verify(s.Ctx(), util.HexAddress{}, []byte{}, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("collections: not found: key '0' of type <nil>"))
		Expect(verify).To(BeFalse())
	})
})
