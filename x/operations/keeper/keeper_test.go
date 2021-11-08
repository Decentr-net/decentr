package keeper

import (
	"testing"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/operations/types"
)

func init() {
	config.SetAddressPrefixes()
}

type keeperSet struct {
	keeper          Keeper
	bankKeeper      types.BankKeeper
	tokenKeeper     types.TokenKeeper
	communityKeeper types.CommunityKeeper
}

func setupKeeper(t testing.TB) (keeperSet, sdk.Context) {
	keys := sdk.NewKVStoreKeys(types.StoreKey, authtypes.StoreKey, banktypes.StoreKey, paramstypes.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	ctx, err := GetContext(keys, tkeys)
	require.NoError(t, err)

	set := keeperSet{}

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsKeeper := paramskeeper.NewKeeper(
		cdc,
		codec.NewLegacyAmino(),
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	accountKeeper := authkeeper.NewAccountKeeper(
		authtypes.ModuleCdc,
		keys[authtypes.StoreKey],
		paramsKeeper.Subspace(authtypes.StoreKey),
		authtypes.ProtoBaseAccount,
		map[string][]string{
			types.ModuleName: {authtypes.Minter, authtypes.Burner},
		},
	)
	set.bankKeeper = bankkeeper.NewBaseKeeper(
		cdc,
		keys[banktypes.StoreKey],
		accountKeeper,
		paramsKeeper.Subspace(banktypes.StoreKey),
		map[string]bool{},
	)

	set.keeper = *NewKeeper(
		cdc,
		keys[types.StoreKey],
		paramsKeeper.Subspace(types.StoreKey),
	)
	set.keeper.SetParams(ctx, types.DefaultParams())

	return set, ctx
}

func TestKeeper_SetParams(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	tt := []struct {
		name    string
		params  types.Params
		isValid bool
	}{
		{
			"default",
			types.DefaultParams(),
			true,
		},
		{
			"zero",
			types.Params{},
			false,
		},
		{
			"invalid_min_gas_price_negative",
			types.Params{
				MinGasPrice: sdk.DecCoin{Denom: config.DefaultBondDenom, Amount: sdk.NewDec(-1)},
			},
			false,
		},
		{
			"invalid_min_gas_price_zero",
			types.Params{
				MinGasPrice: sdk.DecCoin{Denom: config.DefaultBondDenom, Amount: sdk.NewDec(0)},
			},
			false,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.isValid {
				require.NotPanics(t, func() { k.SetParams(ctx, tc.params) })
			} else {
				require.Panics(t, func() { k.SetParams(ctx, tc.params) })
			}
		})
	}
}

func TestKeeper_GetParams(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	require.Equal(t, types.DefaultParams(), k.GetParams(ctx))

	p := types.Params{
		Supervisors: []string{NewAccAddress().String()},
		FixedGas: types.FixedGasParams{
			ResetAccount:      1,
			BanAccount:        2,
			DistributeRewards: 3,
		},
		MinGasPrice: sdk.NewDecCoin(config.DefaultBondDenom, sdk.NewInt(10)),
	}

	k.SetParams(ctx, p)
	require.Equal(t, p, k.GetParams(ctx))
}

func TestKeeper_IsSupervisor(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	a1, a2, a3 := NewAccAddress(), NewAccAddress(), NewAccAddress()

	p := types.DefaultParams()
	p.Supervisors = []string{a1.String(), a2.String()}
	k.SetParams(ctx, p)

	require.True(t, k.IsSupervisor(ctx, a1))
	require.True(t, k.IsSupervisor(ctx, a2))
	require.False(t, k.IsSupervisor(ctx, a3))
}
