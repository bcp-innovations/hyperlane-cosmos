package keeper_test

// Expect(mailbox.Tree.Count).To(Equal(messageSent)) // TODO fix

// TODO fix
//if messageSent == 0 {
//	_, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).LatestCheckpoint(s.Ctx(), &types.QueryLatestCheckpointRequest{Id: mailboxId.String()})
//	Expect(err.Error()).To(Equal("no leaf inserted yet"))
//} else {
//	latestCheckpoint, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).LatestCheckpoint(s.Ctx(), &types.QueryLatestCheckpointRequest{Id: mailboxId.String()})
//	Expect(err).To(BeNil())
//
//	Expect(latestCheckpoint.Count).To(Equal(messageSent - 1))
//}

// TODO: Check claimable fees of IGP
// TODO: check merkle tree hook
