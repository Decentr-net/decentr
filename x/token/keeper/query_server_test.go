package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/token/types"
)

func TestQueryServer_Balance(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewQueryServer(set.keeper, set.distrKeeper)

	address := NewAccAddress()
	set.keeper.IncTokens(ctx, address, sdk.OneDec())

	out, err := s.Balance(sdk.WrapSDKContext(ctx), &types.BalanceRequest{
		Address: address,
	})
	require.NoError(t, err)
	require.Equal(t, &types.BalanceResponse{
		Balance:      sdk.DecProto{sdk.NewDec(2)},
		BalanceDelta: sdk.DecProto{sdk.OneDec()},
		IsBanned:     false,
	}, out)
}

func TestQueryServer_Balance_Banned(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewQueryServer(set.keeper, set.distrKeeper)

	address := NewAccAddress()
	set.keeper.SetBan(ctx, address, true)

	out, err := s.Balance(sdk.WrapSDKContext(ctx), &types.BalanceRequest{
		Address: address,
	})
	require.NoError(t, err)
	require.Equal(t, &types.BalanceResponse{
		Balance:      sdk.DecProto{sdk.OneDec()},
		BalanceDelta: sdk.DecProto{sdk.ZeroDec()},
		IsBanned:     true,
	}, out)
}

func TestQueryServer_Pool(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewQueryServer(set.keeper, set.distrKeeper)

	set.keeper.IncTokens(ctx, NewAccAddress(), sdk.OneDec())

	out, err := s.Pool(sdk.WrapSDKContext(ctx), &types.PoolRequest{})
	require.NoError(t, err)

	require.Equal(t, &types.PoolResponse{
		Size_:                  sdk.NewDecCoin(config.DefaultBondDenom, sdk.NewInt(1000)),
		TotalDelta:             sdk.DecProto{sdk.OneDec()},
		NextDistributionHeight: types.DefaultRewardsBlockInterval,
	}, out)
}
