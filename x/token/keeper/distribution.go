package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/config"
	"github.com/Decentr-net/decentr/x/token/types"
)

func DistributeRewards(ctx sdk.Context, k Keeper, distributionKeeper types.DistributionKeeper) {
	// drop negative deltas
	k.IterateBalanceDelta(ctx, func(address sdk.AccAddress, delta sdk.Dec) (stop bool) {
		if delta.IsNegative() {
			k.SetBalanceDelta(ctx, address, sdk.ZeroDec())
		}
		return false
	})

	total := k.GetAccumulatedDelta(ctx)
	if total.IsZero() {
		return
	}

	pool := distributionKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(config.DefaultBondDenom)
	if pool.IsNil() || pool.IsZero() {
		return
	}

	rate := pool.Quo(total)

	k.IterateBalanceDelta(ctx, func(address sdk.AccAddress, delta sdk.Dec) bool {
		k.SetBalanceDelta(ctx, address, sdk.ZeroDec())

		coins := sdk.Coins{
			sdk.NewCoin(config.DefaultBondDenom, delta.Mul(rate).TruncateInt()),
		}

		if err := distributionKeeper.DistributeFromFeePool(ctx, coins, address); err != nil {
			panic(fmt.Errorf("failed to distribute %s from community pool to %s: %w", coins, address, err))
		}

		if err := ctx.EventManager().EmitTypedEvent(&types.EventRewardDistribution{
			Address: address,
			Delta:   sdk.DecProto{Dec: delta},
			Reward:  coins[0],
		}); err != nil {
			panic(fmt.Errorf("failed to emit event for %s: %w", address, err))
		}

		return false
	})
}
