package types

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	ismTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	pdTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		IsmGenesis:           ismTypes.NewGenesisState(),
		PostDispatchGenesis:  pdTypes.NewGenesisState(),
		Mailboxes:            []Mailbox{},
		Messages:             []GenesisMailboxMessageWrapper{},
		IsmSequence:          0,
		PostDispatchSequence: 0,
		AppSequence:          0,
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	if err := gs.PostDispatchGenesis.Validate(); err != nil {
		return err
	}

	if err := gs.IsmGenesis.Validate(); err != nil {
		return err
	}

	messages := make(map[uint64]map[util.HexAddress]struct{})
	for _, m := range gs.Messages {
		if _, ok := messages[m.MailboxId][m.MessageId]; ok {
			return fmt.Errorf("duplicated message (%s) for mailbox %d, ", m.MessageId, m.MailboxId)
		}
		messages[m.MailboxId][m.MessageId] = struct{}{}
	}

	for i, mailbox := range gs.Mailboxes {
		if mailbox.Id.GetInternalId() != uint64(i) {
			return fmt.Errorf("duplicated mailbox id %d, %d", mailbox.Id.GetInternalId(), i)
		}
		receivedMessages := uint32(len(messages[mailbox.Id.GetInternalId()]))
		if mailbox.MessageReceived != receivedMessages {
			return fmt.Errorf("received messages (%d) != received count (%d) for mailbox %s", receivedMessages, mailbox.MessageReceived, mailbox.Id)
		}
	}

	return nil
}
