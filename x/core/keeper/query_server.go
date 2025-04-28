package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"

	storetypes "cosmossdk.io/store/types"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

func (qs queryServer) Delivered(ctx context.Context, req *types.QueryDeliveredRequest) (*types.QueryDeliveredResponse, error) {
	messageId, err := util.DecodeEthHex(req.MessageId)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	delivered, err := qs.k.Messages.Has(ctx, collections.Join(mailboxId.GetInternalId(), messageId))
	if err != nil {
		return nil, err
	}

	return &types.QueryDeliveredResponse{Delivered: delivered}, nil
}

func (qs queryServer) RecipientIsm(ctx context.Context, req *types.QueryRecipientIsmRequest) (*types.QueryRecipientIsmResponse, error) {
	recipient, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, err
	}

	get, err := qs.k.ReceiverIsmId(ctx, recipient)
	if err != nil {
		return nil, err
	}

	return &types.QueryRecipientIsmResponse{IsmId: get.String()}, nil
}

func (qs queryServer) Mailboxes(ctx context.Context, req *types.QueryMailboxesRequest) (*types.QueryMailboxesResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.Mailboxes, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxesResponse{
		Mailboxes:  values,
		Pagination: pagination,
	}, nil
}

func (qs queryServer) Mailbox(ctx context.Context, req *types.QueryMailboxRequest) (*types.QueryMailboxResponse, error) {
	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

	return &types.QueryMailboxResponse{
		Mailbox: mailbox,
	}, nil
}

func (qs queryServer) VerifyDryRun(ctx context.Context, req *types.QueryVerifyDryRunRequest) (*types.QueryVerifyDryRunResponse, error) {
	limit := uint64(200000)
	if req.GasLimit != "" {
		parsed, err := strconv.ParseUint(req.GasLimit, 10, 32)
		if err != nil {
			return nil, err
		}

		limit = parsed
	}

	// explicitly set a GasMeter to not run into stack overflows as ISM might call themselves
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithGasMeter(storetypes.NewGasMeter(limit))

	ismId, err := util.DecodeHexAddress(req.IsmId)
	if err != nil {
		return nil, err
	}

	// metadata := []byte(req.Metadata)

	bytes, _ := hex.DecodeString("030004671a000000380000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd6d696c6b726f757465725f61707000000000000000000000000000010000000000000000000000000000000000000000894b4cfb1f7f23b864865244b17dbabf333debcb0000000000000000000000000000000000000000000000000000000000002710")
	metadata, _ := hex.DecodeString("000000000000000000000000fdb9cd5f9daaa2e4474019405a328a88e7484f260004671a9e9a8cc0a5a1e41862d63f8105412192cf61ae753517da697c1db2947851c87ca20cfcac9c7d0f1ffa0fe29f33179af70605831e8ad0a2209be326fefa19b9bd582430f5ff154220d5436a922c421ab06286431b315e761057632717f7a6b72db77419caeb42cd9cdb3588453175a2f3d50f6f73bdd17e5b492c06f6d2e8eb815e79fd70b2fc89382c06104a18a850bf5dc9a05b9234888ef0b85ca93d7dd2f3ae518ddee040146c537f148bfb237e56ffdffc3224619a8b93af61ad501c442f998f5d0f0dbe0e1cfc4a733bb241aecd99aacdbb8785db0138499dee1702049e9cdc71a8f2ee6e5c163d01faa020fe8bf6c6f45b45f084b8f8847af2cc6b2cc5ffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f838a7f2483eea821a0d6f2dc8aecb753cc09b1dc111c0f9556f22e9d2ca96d06959c9b639c85b8127a77872d1750c5bc96d01fda35ace2c253f77b67cf8b759403ff474bb281af724ec693f3a3ef829d4041918015537fe5aefc77204166b33398f8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf8923490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99cb3fe49be726b4b5282223ac36460d64cb0a2f7413ed518bb793ac7e03d6b86ea76e183795ad11c544b1047a313c620b5205e2eca9e1ad2588776c79fc8fc93a6da7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d22733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981fe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7ac1452709598b83ec2b101c572a88e698541b722c24a739a3e8392f470eb17443b46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0c65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2f4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd95a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e3774df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652cdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618db8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea32293237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d7358448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a90004674ca57cad376774337b3317878173f911cb62e0488c157262a55fa51f041801712350993d29673b5389f459fa0608b5c2b64555c8a565add7212e93194f45ef34f31cf9883e11c061bc2e1df2455ffd367c3ed9abef58243455ac7496bd6640ed0b5a46e66dc24cf0b69dc1a37fddb244df1ab58dc35b3ef66f00be103f680abf5f541c29edada5d58505d9e52f431a510414a92f32dd78bf8ea70661150fb263d4b17304aaac6677c710db756d8ab00a2fc2b360ac42461e39facefa9412bb7ab08a051c125f4aedea8218f26d29876da2801dad5b17fc8cbc2c34a4528f1a2fc43395005c1d15eff1e33120508d230952992cc1f57d2054c3868db749c24fcef5d681291c")

	msg, err := util.ParseHyperlaneMessage(bytes)
	if err != nil {
		return nil, err
	}

	verified, err := qs.k.Verify(sdkCtx, ismId, metadata, msg)
	return &types.QueryVerifyDryRunResponse{
		Verified: verified,
	}, err
}

func (qs queryServer) RegisteredISMs(_ context.Context, _ *types.QueryRegisteredISMs) (*types.QueryRegisteredISMsResponse, error) {
	return &types.QueryRegisteredISMsResponse{
		Ids: qs.k.IsmRouter().GetModuleIds(),
	}, nil
}

func (qs queryServer) RegisteredHooks(_ context.Context, _ *types.QueryRegisteredHooks) (*types.QueryRegisteredHooksResponse, error) {
	return &types.QueryRegisteredHooksResponse{
		Ids: qs.k.PostDispatchRouter().GetModuleIds(),
	}, nil
}

func (qs queryServer) RegisteredApps(_ context.Context, _ *types.QueryRegisteredApps) (*types.QueryRegisteredAppsResponse, error) {
	return &types.QueryRegisteredAppsResponse{
		Ids: qs.k.AppRouter().GetModuleIds(),
	}, nil
}
