package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"

	"github.com/Decentr-net/decentr/x/community/types"
)

func getPostKeeperKey(owner sdk.AccAddress, id uuid.UUID) []byte {
	return append(owner.Bytes(), id.Bytes()[:]...)
}

// CreatePost creates new post. Keeper's key is joined owner and uuid.
func (k Keeper) CreatePost(ctx sdk.Context, p types.Post) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	key := getPostKeeperKey(p.Owner, p.UUID)
	store.Set(key, k.cdc.MustMarshalBinaryBare(p))
}

// DeletePost deletes the post from keeper.
func (k Keeper) DeletePost(ctx sdk.Context, owner sdk.AccAddress, id uuid.UUID) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	key := append(owner.Bytes(), id.Bytes()...)

	likesStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)
	it := sdk.KVStorePrefixIterator(likesStore, key)
	for ; it.Valid(); it.Next() {
		likesStore.Delete(it.Key())
	}
	it.Close()

	store.Delete(key)
}

// GetPostByKey returns entire post by keeper's key.
func (k Keeper) GetPostByKey(ctx sdk.Context, key []byte) types.Post {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	if !store.Has(key) {
		return types.Post{}
	}

	bz := store.Get(key)

	var post types.Post
	k.cdc.MustUnmarshalBinaryBare(bz, &post)
	return post
}

// GetPostsIterator returns an iterator over all posts
func (k Keeper) GetPostsIterator(ctx sdk.Context) sdk.Iterator {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix).Iterator(nil, nil)
}
