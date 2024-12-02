package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/KYVENetwork/mailbox"
	"github.com/KYVENetwork/mailbox/keeper"
)

type testFixture struct {
	ctx         sdk.Context
	k           keeper.Keeper
	msgServer   mailbox.MsgServer
	queryServer mailbox.QueryServer

	addrs []sdk.AccAddress
}

func initFixture(t *testing.T) *testFixture {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	key := storetypes.NewKVStoreKey(mailbox.ModuleName)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	storeService := runtime.NewKVStoreService(key)
	addrs := simtestutil.CreateIncrementalAccounts(3)

	k := keeper.NewKeeper(encCfg.Codec, addresscodec.NewBech32Codec("cosmos"), storeService, addrs[0].String())
	err := k.InitGenesis(testCtx.Ctx, mailbox.NewGenesisState())
	if err != nil {
		panic(err)
	}

	return &testFixture{
		ctx:         testCtx.Ctx,
		k:           k,
		msgServer:   keeper.NewMsgServerImpl(k),
		queryServer: keeper.NewQueryServerImpl(k),
		addrs:       addrs,
	}
}
