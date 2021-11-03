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
