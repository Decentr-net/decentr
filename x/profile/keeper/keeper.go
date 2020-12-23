package keeper

import (
	"github.com/Decentr-net/decentr/x/profile/types"
	"github.com/Decentr-net/decentr/x/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int, description []byte)
}

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	tokens   TokenKeeper
}

// NewKeeper creates new instances of the PDV Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, tokens TokenKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		tokens:   tokens,
		storeKey: storeKey,
	}
}

// Sets the entire Profile metadata struct for a name
func (k Keeper) SetProfile(ctx sdk.Context, owner sdk.AccAddress, new types.Profile) {
	store := ctx.KVStore(k.storeKey)

	old := k.GetProfile(ctx, owner)
	if old.Public.RegisteredAt == 0 {
		new.Public.RegisteredAt = ctx.BlockTime().Unix()
		k.tokens.AddTokens(ctx, owner, utils.InitialTokenBalance(), []byte("profile"))
	} else {
		new.Public.RegisteredAt = old.Public.RegisteredAt
	}

	store.Set(owner, k.cdc.MustMarshalBinaryBare(new))
}

// Gets the entire Profile metadata struct for an owner
func (k Keeper) GetProfile(ctx sdk.Context, owner sdk.AccAddress) types.Profile {
	store := ctx.KVStore(k.storeKey)

	if !k.IsNamePresent(ctx, owner) {
		return types.Profile{Owner: owner}
	}

	bz := store.Get(owner)

	var profile types.Profile
	k.cdc.MustUnmarshalBinaryBare(bz, &profile)
	return profile
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, owner sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(owner)
}

// Get an iterator over all names in which the keys are the names and the values are the Profile
func (k Keeper) GetProfileIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
