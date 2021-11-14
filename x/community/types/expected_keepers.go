package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type TokenKeeper interface {
	IncTokens(ctx sdk.Context, address sdk.AccAddress, amount sdk.Dec)
}
