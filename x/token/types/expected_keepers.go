package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DistributionKeeper interface {
	GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins
	DistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, address sdk.AccAddress) error
}
