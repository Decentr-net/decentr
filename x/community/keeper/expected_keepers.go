package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int)
	GetBalance(ctx sdk.Context, owner sdk.AccAddress) sdk.Int
}
