package types

import (
	"context"
	//"fmt"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// combine multiple mailbox hooks, all hook functions are run in array sequence
var _ MailboxHooks = MultiMailboxHooks{}

type MultiMailboxHooks []MailboxHooks

func NewMultiMailboxHooks(hooks ...MailboxHooks) MultiMailboxHooks {
	return hooks
}

func (h MultiMailboxHooks) Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message util.HyperlaneMessage) error {
	for i := range h {
		if err := h[i].Handle(ctx, mailboxId, origin, sender, message); err != nil {
			return err
		}
	}

	return nil
}

//hook.PostDispatch()
//
//keeper.PostDispatch(id) {
//	hook <- store(id)
//
//
//	keeper.collection.Save(hook)
//}
//core()
//hook.PostDispatch() updatedHook,

//type HyperlanePostDispatchHook interface {
//	proto.Message
//
//	PostDispatch(metadata any, message HyperlaneMessage) error
//	QuoteDispatch(metadata any, message HyperlaneMessage) (error, sdk.Coins)
//}
//
//func (m *MerkleTreeHook) PostDispatch(options any, message HyperlaneMessage) error {
//	f, ok := options.(MerkleTreeHookOptions)
//	if !ok {
//		return fmt.Errorf("options should be MerkleTreeHookOptions")
//	}
//	_ = f.Test
//
//	return nil
//}
//
//func (m *MerkleTreeHook) QuoteDispatch(options any, message HyperlaneMessage) (error, sdk.Coins) {
//	f, ok := options.(MerkleTreeHookOptions)
//	if !ok {
//		return fmt.Errorf("options should be MerkleTreeHookOptions"), nil
//	}
//	_ = f.Test
//
//	return nil, nil
//}
//
//type MerkleTreeHookOptions struct {
//	Test string
//}
