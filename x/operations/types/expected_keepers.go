package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type TokenKeeper interface {
	IncTokens(ctx sdk.Context, address sdk.AccAddress, amount sdk.Dec)
	ResetAccount(ctx sdk.Context, address sdk.AccAddress)
}

type CommunityKeeper interface {
	ResetAccount(ctx sdk.Context, address sdk.AccAddress)
}
