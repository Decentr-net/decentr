package keeper

import (
	"fmt"

	"github.com/Decentr-net/decentr/config"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Decentr-net/decentr/x/token/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	paramSpace paramtypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) IncTokens(ctx sdk.Context, address sdk.AccAddress, amount sdk.Dec) {
	if amount.IsZero() {
		return
	}

	k.SetBalance(ctx, address, k.GetBalance(ctx, address).Add(amount))
	k.SetBalanceDelta(ctx, address, k.GetBalanceDelta(ctx, address).Add(amount))
}

func (k Keeper) GetBalance(ctx sdk.Context, address sdk.AccAddress) sdk.Dec {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalancePrefix)

	if !store.Has(address) {
		return config.InitialTokenBalance
	}

	var d sdk.Dec
	if err := d.Unmarshal(store.Get(address)); err != nil {
		panic(fmt.Errorf("failed to get balance for %s: %w", address, err))
	}
	return d
}

func (k Keeper) SetBalance(ctx sdk.Context, address sdk.AccAddress, amount sdk.Dec) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalancePrefix)

	if amount.IsZero() {
		store.Delete(address)
		return
	}

	bz, err := amount.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal new balance for %s: %w", address, err))
	}
	store.Set(address, bz)
}

func (k Keeper) GetAccumulatedDelta(ctx sdk.Context) sdk.Dec {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CalcPrefix)

	if !store.Has(types.AccumulatedDeltaKey) {
		return sdk.ZeroDec()
	}

	var d sdk.Dec
	if err := d.Unmarshal(store.Get(types.AccumulatedDeltaKey)); err != nil {
		panic(fmt.Errorf("failed to get accumulated delta: %w", err))
	}
	return d
}

func (k Keeper) IncAccumulatedDelta(ctx sdk.Context, amount sdk.Dec) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CalcPrefix)

	bz, err := k.GetAccumulatedDelta(ctx).Add(amount).Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal new accumulated delta: %w", err))
	}
	store.Set(types.AccumulatedDeltaKey, bz)
}

func (k Keeper) SetBalanceDelta(ctx sdk.Context, address sdk.AccAddress, amount sdk.Dec) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalanceDeltaPrefix)

	if diff := amount.Sub(k.GetBalanceDelta(ctx, address)); !diff.IsZero() {
		k.IncAccumulatedDelta(ctx, diff)
	}

	if amount.IsZero() {
		store.Delete(address)
		return
	}

	bz, err := amount.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal delta amount for %s: %w", address, amount))
	}
	store.Set(address, bz)
}

func (k Keeper) GetBalanceDelta(ctx sdk.Context, address sdk.AccAddress) sdk.Dec {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalanceDeltaPrefix)

	if !store.Has(address) {
		return sdk.ZeroDec()
	}

	var d sdk.Dec
	if err := d.Unmarshal(store.Get(address)); err != nil {
		panic(fmt.Errorf("failed to get balance for %s: %w", address, err))
	}
	return d
}

func (k Keeper) ResetAccount(ctx sdk.Context, address sdk.AccAddress) {
	k.SetBalanceDelta(ctx, address, sdk.ZeroDec())
	k.SetBalance(ctx, address, sdk.ZeroDec())
}

func (k Keeper) SetBan(ctx sdk.Context, address sdk.AccAddress, ban bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BanListPrefix)

	if ban {
		k.SetBalanceDelta(ctx, address, sdk.ZeroDec())
		store.Set(address, []byte{0x00})
	} else {
		store.Delete(address)
	}
}

func (k Keeper) IsBanned(ctx sdk.Context, address sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BanListPrefix)

	return store.Has(address)
}

func (k Keeper) IterateBalance(ctx sdk.Context, handle func(address sdk.AccAddress, balance sdk.Dec) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalancePrefix)

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		address := sdk.AccAddress(it.Key())

		var delta sdk.Dec
		if err := delta.Unmarshal(it.Value()); err != nil {
			panic(fmt.Errorf("failed to get balance for %s: %w", address, err))
		}

		if handle(it.Key(), delta) {
			return
		}

		it.Next()
	}
}

func (k Keeper) IterateBalanceDelta(ctx sdk.Context, handle func(address sdk.AccAddress, delta sdk.Dec) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalanceDeltaPrefix)

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		address := sdk.AccAddress(it.Key())

		var delta sdk.Dec
		if err := delta.Unmarshal(it.Value()); err != nil {
			panic(fmt.Errorf("failed to get delta for %s: %w", address, err))
		}

		if handle(it.Key(), delta) {
			return
		}

		it.Next()
	}
}

func (k Keeper) IterateBanList(ctx sdk.Context, handle func(address sdk.AccAddress) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BanListPrefix)

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		if handle(it.Key()) {
			return
		}

		it.Next()
	}
}
