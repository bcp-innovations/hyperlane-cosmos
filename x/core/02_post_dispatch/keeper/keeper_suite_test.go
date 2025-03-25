package keeper_test

import (
	"fmt"
	"testing"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	ismTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPostDispatchKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.SubModuleName))
}

func createNoopISM(s *i.KeeperTestSuite, creator string) (util.HexAddress, error) {
	res, err := s.RunTx(&ismTypes.MsgCreateNoopIsm{Creator: creator})
	if err != nil {
		return [32]byte{}, err
	}

	var response ismTypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	if err != nil {
		return [32]byte{}, err
	}

	return response.Id, nil
}

func createDummyMailbox(s *i.KeeperTestSuite, creator string) (util.HexAddress, error) {
	ismId, err := createNoopISM(s, creator)
	if err != nil {
		return [32]byte{}, err
	}

	res, err := s.RunTx(&coreTypes.MsgCreateMailbox{
		Owner:        creator,
		LocalDomain:  11,
		DefaultIsm:   ismId,
		DefaultHook:  nil,
		RequiredHook: nil,
	})
	if err != nil {
		return [32]byte{}, err
	}

	var response coreTypes.MsgCreateMailboxResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	if err != nil {
		return [32]byte{}, err
	}

	mailboxId, err := util.DecodeHexAddress(response.Id)
	if err != nil {
		return [32]byte{}, err
	}

	return mailboxId, nil
}

func createDummyMerkleTreeHook(s *i.KeeperTestSuite, creator string, mailboxId util.HexAddress) (util.HexAddress, error) {
	res, err := s.RunTx(&types.MsgCreateMerkleTreeHook{
		Owner:     creator,
		MailboxId: mailboxId.String(),
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateMerkleTreeHookResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	if err != nil {
		return [32]byte{}, err
	}

	hookId, err := util.DecodeHexAddress(response.Id)
	if err != nil {
		return [32]byte{}, err
	}

	return hookId, nil
}
