package integration

import (
	"github.com/bcp-innovations/hyperlane-cosmos/tests/simapp"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"time"
)

var (
	A_DENOM = "acoin"
	B_DENOM = "bcoin"
	C_DENOM = "ccoin"
)

func NewCleanChain() *KeeperTestSuite {
	return NewCleanChainAtTime(time.Now().Unix())
}

func NewCleanChainAtTime(startTime int64) *KeeperTestSuite {
	s := KeeperTestSuite{}
	s.setupApp(startTime)
	return &s
}

func (suite *KeeperTestSuite) App() *simapp.App {
	return suite.app
}

func (suite *KeeperTestSuite) Ctx() sdk.Context {
	return suite.ctx
}

func (suite *KeeperTestSuite) Commit() {
	suite.commitAfter(time.Second * 0)
}

func (suite *KeeperTestSuite) CommitAfterSeconds(seconds uint64) {
	suite.commitAfter(time.Second * time.Duration(seconds))
}

func (suite *KeeperTestSuite) commitAfter(t time.Duration) {
	header := suite.ctx.BlockHeader()
	header.Time = header.Time.Add(t)

	_, err := suite.app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: header.Height,
		Time:   header.Time,
		DecidedLastCommit: abci.CommitInfo{
			Round: 0,
			Votes: suite.VoteInfos,
		},
	})
	if err != nil {
		panic(err)
	}
	_, err = suite.app.Commit()
	if err != nil {
		panic(err)
	}

	header.Height += 1

	suite.ctx = suite.app.BaseApp.NewUncachedContext(false, header)
}

// ##########################
// #          MINT          #
// ##########################

type TestValidatorAddress struct {
	Moniker string

	PrivateKey *ed25519.PrivKey

	Address        string
	AccAddress     sdk.AccAddress
	ConsAccAddress sdk.ConsAddress
	ConsAddress    string
}

func (suite *KeeperTestSuite) MintBaseCoins(address string, amount uint64) error {
	return suite.MintCoins(address, sdk.NewCoins(
		// mint coins A, B, C
		sdk.NewInt64Coin(A_DENOM, int64(amount)),
		sdk.NewInt64Coin(B_DENOM, int64(amount)),
		sdk.NewInt64Coin(C_DENOM, int64(amount)),
	))
}

func (suite *KeeperTestSuite) MintCoins(address string, coins sdk.Coins) error {
	// mint coins A, B, C
	err := suite.app.BankKeeper.MintCoins(suite.ctx, bankTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	receiver, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, bankTypes.ModuleName, receiver, coins)
	if err != nil {
		return err
	}

	return nil
}

func GenerateTestValidatorAddress(moniker string) TestValidatorAddress {
	a := TestValidatorAddress{}
	a.Moniker = moniker
	a.PrivateKey = ed25519.GenPrivKeyFromSecret([]byte(moniker))

	a.AccAddress = sdk.AccAddress(a.PrivateKey.PubKey().Address())
	bech32Address, _ := sdk.Bech32ifyAddressBytes(simapp.Bech32PrefixAccAddr, a.AccAddress)
	a.Address = bech32Address

	a.ConsAccAddress = sdk.ConsAddress(a.PrivateKey.PubKey().Address())
	bech32ConsAddress, _ := sdk.Bech32ifyAddressBytes(simapp.Bech32PrefixConsAddr, a.AccAddress)
	a.ConsAddress = bech32ConsAddress

	return a
}
