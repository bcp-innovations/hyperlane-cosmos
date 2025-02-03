package integration

import (
	"encoding/json"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/tests/simapp"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cmtProto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cmtTypes "github.com/cometbft/cometbft/types"
	"github.com/cometbft/cometbft/version"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperTestSuite struct {
	ctx sdk.Context

	app                 *simapp.App
	denom               string
	privateValidatorKey *ed25519.PrivKey
	VoteInfos           []abci.VoteInfo
}

// DefaultConsensusParams ...
var DefaultConsensusParams = &cmtProto.ConsensusParams{
	Block: &cmtProto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1, // no limit
	},
	Evidence: &cmtProto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &cmtProto.ValidatorParams{
		PubKeyTypes: []string{
			cmtTypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func (suite *KeeperTestSuite) setupApp(startTime int64) {
	db := dbm.NewMemDB()

	logger := log.NewNopLogger()
	localApp, err := simapp.NewMiniApp(logger, db, nil, true, EmptyAppOptions{}, baseapp.SetChainID("hyperlane-local"))
	if err != nil {
		panic(err)
	}
	suite.app = localApp

	suite.privateValidatorKey = ed25519.GenPrivKeyFromSecret([]byte("Validator-1"))
	genesisState := DefaultGenesisWithValSet(suite.app, suite.privateValidatorKey)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	if _, err = suite.app.InitChain(
		&abci.RequestInitChain{
			ChainId:         "hyperlane-local",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	); err != nil {
		panic(err)
	}

	suite.denom = simapp.CoinUnit

	suite.ctx = suite.app.BaseApp.NewContextLegacy(false, tmproto.Header{
		Height:          1,
		ChainID:         "hyperlane-local",
		Time:            time.Unix(startTime, 0).UTC(),
		ProposerAddress: sdk.ConsAddress(suite.privateValidatorKey.PubKey().Address()).Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})
}

func DefaultGenesisWithValSet(app *simapp.App, validatorPrivateKey *ed25519.PrivKey) map[string]json.RawMessage {
	bondingDenom := simapp.DefaultBondDenom

	// Generate a new validator.
	pubKey := validatorPrivateKey.PubKey()
	valAddress := sdk.ValAddress(pubKey.Address()).String()
	pkAny, _ := codectypes.NewAnyWithValue(pubKey)

	validators := []stakingTypes.Validator{
		{
			OperatorAddress:   valAddress,
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingTypes.Bonded,
			Tokens:            sdk.DefaultPowerReduction,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingTypes.Description{},
			UnbondingHeight:   0,
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingTypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		},
	}
	// Generate a new delegator.
	delegator := authTypes.NewBaseAccount(
		validatorPrivateKey.PubKey().Address().Bytes(), validatorPrivateKey.PubKey(), 0, 0,
	)

	delegations := []stakingTypes.Delegation{
		stakingTypes.NewDelegation(delegator.GetAddress().String(), valAddress, math.LegacyOneDec()),
	}

	// Default genesis state.
	genesisState := app.DefaultGenesis()

	// Update x/auth state.
	authGenesis := authTypes.NewGenesisState(authTypes.DefaultParams(), []authTypes.GenesisAccount{delegator})
	genesisState[authTypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	// Update x/bank state.
	bondedCoins := sdk.NewCoins(sdk.NewCoin(bondingDenom, sdk.DefaultPowerReduction))

	bankGenesis := bankTypes.NewGenesisState(bankTypes.DefaultGenesisState().Params, []bankTypes.Balance{
		{
			Address: authTypes.NewModuleAddress(stakingTypes.BondedPoolName).String(),
			Coins:   bondedCoins,
		},
	}, bondedCoins, []bankTypes.Metadata{}, []bankTypes.SendEnabled{})
	genesisState[bankTypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Update x/staking state.
	stakingParams := stakingTypes.DefaultParams()
	stakingParams.BondDenom = bondingDenom
	stakingParams.MaxValidators = 1

	stakingGenesis := stakingTypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingTypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	return genesisState
}

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} { return nil }
