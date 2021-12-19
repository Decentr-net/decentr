package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/token/types"
)

func init() {
	config.SetAddressPrefixes()
}

type keeperSet struct {
	keeper     Keeper
	bankKeeper bankkeeper.Keeper
}

func setupKeeper(t testing.TB) (keeperSet, sdk.Context) {
	keys := sdk.NewKVStoreKeys(
		types.StoreKey, authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distrtypes.StoreKey, paramstypes.StoreKey)
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
			distrtypes.ModuleName:          nil,
			minttypes.ModuleName:           {authtypes.Minter},
			stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
			stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
			types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		},
	)

	set.bankKeeper = bankkeeper.NewBaseKeeper(
		cdc,
		keys[banktypes.StoreKey],
		accountKeeper,
		paramsKeeper.Subspace(banktypes.StoreKey),
		map[string]bool{},
	)
	require.NoError(t, set.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000)),
	)))

	stakingKeeper := stakingkeeper.NewKeeper(
		cdc,
		keys[stakingtypes.StoreKey],
		accountKeeper,
		set.bankKeeper,
		paramsKeeper.Subspace(stakingtypes.StoreKey),
	)

	distrKeeper := distrkeeper.NewKeeper(
		cdc,
		keys[distrtypes.StoreKey],
		paramsKeeper.Subspace(distrtypes.StoreKey),
		accountKeeper,
		set.bankKeeper,
		stakingKeeper,
		"collector",
		map[string]bool{},
	)
	distrKeeper.SetFeePool(ctx, distrtypes.FeePool{CommunityPool: nil})
	require.NoError(t, distrKeeper.FundCommunityPool(
		ctx,
		sdk.NewCoins(
			sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000)),
		),
		accountKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress(),
	))

	set.keeper = *NewKeeper(
		cdc,
		keys[types.StoreKey],
	)

	return set, ctx
}

func TestKeeper_SetBalance(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	address := NewAccAddress()

	k.SetBalance(ctx, address, sdk.MustNewDecFromStr("0.000001"))

	require.Equal(t, sdk.NewDecWithPrec(1, 6), k.GetBalance(ctx, address))
}

func TestKeeper_GetBalance(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	// default balance
	require.Equal(t, sdk.OneDec(), k.GetBalance(ctx, NewAccAddress()))
}

func TestKeeper_IterateBalance(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	exp := make(map[string]sdk.Dec, 10)
	for i := 0; i < 10; i++ {
		address := NewAccAddress()
		balance := sdk.NewDec(int64(i + 1))
		exp[address.String()] = balance
		k.SetBalance(ctx, address, balance)
	}

	act := make(map[string]sdk.Dec, 10)
	k.IterateBalance(ctx, func(address sdk.AccAddress, balance sdk.Dec) (stop bool) {
		act[address.String()] = balance
		return
	})

	require.Equal(t, exp, act)
}

func TestKeeper_IncTokens(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	address := NewAccAddress()

	k.IncTokens(ctx, address, sdk.OneDec())

	require.Equal(t, sdk.NewDec(2), k.GetBalance(ctx, address)) // pdv = 1.0 is initial
}

func TestKeeper_ResetAccount(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	address := NewAccAddress()
	k.IncTokens(ctx, address, sdk.OneDec())

	k.ResetAccount(ctx, address)

	require.Equal(t, sdk.OneDec(), k.GetBalance(ctx, address)) // pdv = 1.0 is initial
}
