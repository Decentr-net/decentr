package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/types"
)

func (k Keeper) Follow(ctx sdk.Context, who, whom sdk.Address) {
	if who.Equals(whom) {
		return
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowersPrefix)
	prefix.NewStore(store, who.Bytes()).Set(whom.Bytes(), []byte{})
}

func (k Keeper) Unfollow(ctx sdk.Context, who, whom sdk.Address) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowersPrefix)
	prefix.NewStore(store, who.Bytes()).Delete(whom.Bytes())
}

func (k Keeper) GetFollowees(ctx sdk.Context, who sdk.Address) []sdk.Address {
	var out []sdk.Address

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowersPrefix)
	it := sdk.KVStorePrefixIterator(prefix.NewStore(store, who.Bytes()), nil)
	for ; it.Valid(); it.Next() {
		out = append(out, sdk.AccAddress(it.Key()))
	}
	it.Close()

	return out
}

// IterateFollowers provide iterator over all followers
func (k Keeper) IterateFollowers(ctx sdk.Context, cb func(who, whom sdk.Address) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowersPrefix)
	it := store.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		// it is a bit hacky, it.Key() is composite: consist of concatenated who/whom account keys
		who := sdk.AccAddress(it.Key()[0:sdk.AddrLen])
		whom := sdk.AccAddress(it.Key()[sdk.AddrLen:])
		if cb(who, whom) {
			break
		}
	}
}
