package types

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

type CoreKeeper interface {
	LocalDomain(ctx context.Context) (uint32, error)
	MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error)
	AppRouter() *util.Router[util.HyperlaneApp]
	DispatchMessage(
		ctx sdk.Context,
		originMailboxId util.HexAddress,
		// sender address on the origin chain (e.g. token id)
		sender util.HexAddress,
		// the maximum amount of tokens the dispatch is allowed to cost
		maxFee sdk.Coins,
		destinationDomain uint32,
		// Recipient address on the destination chain (e.g. smart contract)
		recipient util.HexAddress,
		body []byte,
		// Custom metadata for postDispatch Hook
		metadata []byte,
		postDispatchHookId util.HexAddress,
	) (messageId util.HexAddress, error error)
}
