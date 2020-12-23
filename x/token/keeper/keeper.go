package keeper

import (
	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/Decentr-net/decentr/x/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var totalSupplyKey = []byte("totalSupply")

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// AddTokens adds token to the given owner
// Description is needed to merge records in the index.
func (k Keeper) AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int, description []byte) {
	store := ctx.KVStore(k.storeKey)

	balance := k.GetBalance(ctx, owner)
	balance = balance.Add(amount)

	store.Set(getStoreKey(owner), k.cdc.MustMarshalBinaryBare(balance))

	timestamp := uint64(ctx.BlockTime().Unix())
	statsKey := getStatsKey(append(append(owner, utils.Uint64ToBytes(timestamp)...), description...))

	store.Set(statsKey, k.cdc.MustMarshalBinaryBare(amount))

	k.addTotalSupply(ctx, amount)
}

// addTotalSupply increase or decrease total supply with the given amount of tokens
func (k Keeper) addTotalSupply(ctx sdk.Context, amount sdk.Int) {
	balance := k.GetTotalSupply(ctx)
	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(totalSupplyKey, k.cdc.MustMarshalBinaryBare(balance))
}

// GetBalance returns token balance for the given owner
func (k Keeper) GetBalance(ctx sdk.Context, owner sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	if store.Has(getStoreKey(owner)) {
		balance := store.Get(getStoreKey(owner))
		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(balance, &amount)
		return amount
	}
	return sdk.ZeroInt()
}

// GetTotalSupply returns total token supply
func (k Keeper) GetTotalSupply(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	totalSupply := store.Get(totalSupplyKey)
	if totalSupply == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.Int
	k.cdc.MustUnmarshalBinaryBare(totalSupply, &amount)
	return amount
}

// Get an iterator over all balances in which the keys are the accounts and the values are their balance
func (k Keeper) GetBalanceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.StorePrefix)
}

func (k Keeper) GetStats(ctx sdk.Context, owner sdk.AccAddress) map[uint64]float64 {
	store := ctx.KVStore(k.storeKey)

	it := sdk.KVStorePrefixIterator(store, getStatsKey(owner))

	t := uint64(0)
	a := sdk.NewInt(0)

	out := make(map[uint64]float64, 365)
	for ; it.Valid(); it.Next() {
		timestamp := utils.BytesToUint64(it.Key()[:8]) // first part of key is timestamp, second - random uuid
		truncated := timestamp - timestamp%86400       // truncate to day

		if t == 0 {
			t = truncated
		}

		if t != truncated {
			out[t] = utils.TokenToFloat64(a)
			t = truncated
		}

		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &amount)

		a = a.Add(amount)
	}

	return out
}

func getStoreKey(key []byte) []byte {
	return append(types.StorePrefix, key...)
}

func getStatsKey(key []byte) []byte {
	return append(types.StatsPrefix, key...)
}
