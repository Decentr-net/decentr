package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	AddTokens(ctx sdk.Context, owner sdk.AccAddress, timestamp time.Time, amount sdk.Int)
}
