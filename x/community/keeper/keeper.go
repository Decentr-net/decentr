package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	tokens   TokenKeeper

	paramSubspace params.Subspace
}

// NewKeeper creates new instances of the community Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace, tokens TokenKeeper) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		tokens:        tokens,
		paramSubspace: ps,
	}
}

func (k *Keeper) ResetAccount(ctx sdk.Context, owner sdk.AccAddress) {
	k.IterateFollowers(ctx, func(who, whom sdk.Address) (stop bool) {
		if who.Equals(owner) || whom.Equals(owner) {
			k.Unfollow(ctx, who, whom)
		}

		return false
	})

	it := k.GetPostsIterator(ctx)
	for ; it.Valid(); it.Next() {
		if p := k.GetPostByKey(ctx, it.Key()); p.Owner.Equals(owner) {
			k.DeletePost(ctx, owner, p.UUID)
		}
	}
	it.Close()

	it = k.GetLikesIterator(ctx)
	for ; it.Valid(); it.Next() {
		l := k.GetLikeByKey(ctx, it.Key())

		if l.PostOwner.Equals(owner) || l.Owner.Equals(owner) {
			k.DeleteLike(ctx, l)
		}
	}
}
