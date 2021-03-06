package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/types"
)

func getLikeKeeperKey(l types.Like) []byte {
	return append(append(l.PostOwner.Bytes(), l.PostUUID[:]...), l.Owner.Bytes()...)
}

func (k Keeper) SetLike(ctx sdk.Context, l types.Like) {
	if l.Owner.Equals(l.PostOwner) {
		ctx.Logger().Info("SetLike: owner tries to like own post",
			"postUUID", l.PostUUID,
			"postOwner", l.PostOwner)
		return
	}

	likesStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)

	key := getLikeKeeperKey(l)
	old := k.GetLikeByKey(ctx, key)

	k.tokens.AddTokens(ctx, l.PostOwner, sdk.NewInt(int64(l.Weight)-int64(old.Weight)))

	likesStore.Set(key, k.cdc.MustMarshalBinaryBare(l))
}

// GetLikeByKey returns entire like by keeper's key.
func (k Keeper) GetLikeByKey(ctx sdk.Context, key []byte) types.Like {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)

	if !store.Has(key) {
		return types.Like{}
	}

	bz := store.Get(key)

	var like types.Like
	k.cdc.MustUnmarshalBinaryBare(bz, &like)
	return like
}

// GetLikesIterator returns an iterator over all likes
func (k Keeper) GetLikesIterator(ctx sdk.Context) sdk.Iterator {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix).Iterator(nil, nil)
}
