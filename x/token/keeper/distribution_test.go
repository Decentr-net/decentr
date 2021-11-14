package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
)

func TestDistributeRewards(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	addr1, addr2 := NewAccAddress(), NewAccAddress()
	k.IncTokens(ctx, addr1, sdk.OneDec())
	k.IncTokens(ctx, addr2, sdk.OneDec())

	DistributeRewards(ctx, k, set.distrKeeper)
	require.Empty(t, set.distrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, sdk.ZeroDec(), k.GetAccumulatedDelta(ctx))
	require.Equal(t, sdk.ZeroDec(), k.GetBalanceDelta(ctx, addr1))
	require.Equal(t, sdk.ZeroDec(), k.GetBalanceDelta(ctx, addr2))

	require.Equal(t, sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(500)), set.bankKeeper.GetBalance(ctx, addr1, config.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(500)), set.bankKeeper.GetBalance(ctx, addr2, config.DefaultBondDenom))
}

func TestDistributeRewards_NegativeDelta(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	addr1, addr2 := NewAccAddress(), NewAccAddress()

	// that account shouldn't affect airdrop
	k.IncTokens(ctx, NewAccAddress(), sdk.OneDec().Neg())

	// result of code below should be the same with TestDistributeRewards
	k.IncTokens(ctx, addr1, sdk.OneDec())
	k.IncTokens(ctx, addr2, sdk.OneDec())

	DistributeRewards(ctx, k, set.distrKeeper)
	require.Empty(t, set.distrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, sdk.ZeroDec(), k.GetAccumulatedDelta(ctx))
	require.Equal(t, sdk.ZeroDec(), k.GetBalanceDelta(ctx, addr1))
	require.Equal(t, sdk.ZeroDec(), k.GetBalanceDelta(ctx, addr2))

	require.Equal(t, sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(500)), set.bankKeeper.GetBalance(ctx, addr1, config.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(500)), set.bankKeeper.GetBalance(ctx, addr2, config.DefaultBondDenom))
}

func TestDistributeRewards_ZeroAccumulated(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	k.IncTokens(ctx, NewAccAddress(), sdk.OneDec().Neg())
	k.IncTokens(ctx, NewAccAddress(), sdk.OneDec().Neg())

	DistributeRewards(ctx, k, set.distrKeeper)
	require.Equal(t, sdk.DecCoins{sdk.NewDecCoin(config.DefaultBondDenom, sdk.NewInt(1000))}, set.distrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, sdk.ZeroDec(), k.GetAccumulatedDelta(ctx))
}
