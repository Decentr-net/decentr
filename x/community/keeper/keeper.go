package keeper

import (
	"github.com/Decentr-net/decentr/x/community/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"
)

type TokenKeeper interface {
	AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int)
}

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	index    Index
	tokens   TokenKeeper
}

// NewKeeper creates new instances of the community Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, index Index, tokens TokenKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		tokens:   tokens,
		index:    index,
	}
}

// CreatePost creates new post. Keeper's key is joined owner and uuid.
func (k Keeper) CreatePost(ctx sdk.Context, p types.Post) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	k.index.AddPost(p)
	store.Set(getPostKeeperKey(p), k.cdc.MustMarshalBinaryBare(p))
}

// DeletePost deletes the post from keeper.
func (k Keeper) DeletePost(ctx sdk.Context, owner sdk.AccAddress, id uuid.UUID) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	key := append(owner.Bytes(), id.Bytes()...)

	k.index.DeletePost(k.GetPostByKey(ctx, key))
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
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.PostPrefix)
}

func (k Keeper) SetLike(ctx sdk.Context, newLike types.Like) {
	if newLike.Owner.Equals(newLike.PostOwner) {
		// post owner can not like their post
		return
	}

	postsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)
	oldPostBytes := postsStore.Get(getPostKeeperKeyFromLike(newLike))

	if oldPostBytes == nil {
		// post does not exists
		return
	}

	var oldPost types.Post
	k.cdc.MustUnmarshalBinaryBare(oldPostBytes, &oldPost)

	var oldLike types.Like
	oldLike.Weight = types.LikeWeightZero

	likesStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)
	oldLikeBytes := likesStore.Get(getLikeKeeperKey(newLike))
	if oldLikeBytes != nil {
		k.cdc.MustUnmarshalBinaryBare(oldLikeBytes, &oldLike)
	}

	// copy old post to new post
	newPost := oldPost
	updatePostLikesCounters(&newPost, oldLike.Weight, newLike.Weight)

	postsStore.Set(getPostKeeperKey(newPost), k.cdc.MustMarshalBinaryBare(newPost))
	likesStore.Set(getLikeKeeperKey(newLike), k.cdc.MustMarshalBinaryBare(newLike))
	k.index.UpdateLikes(oldPost, newPost)
}

func updatePostLikesCounters(post *types.Post, oldWeight types.LikeWeight, newWeight types.LikeWeight) {
	switch oldWeight {
	case types.LikeWeightUp:
		switch newWeight {
		case types.LikeWeightDown:
			post.LikesCount -= 1
			post.DislikesCount += 1
		case types.LikeWeightZero:
			post.LikesCount -= 1
		}
	case types.LikeWeightDown:
		switch newWeight {
		case types.LikeWeightUp:
			post.LikesCount += 1
			post.DislikesCount -= 1
		case types.LikeWeightZero:
			post.DislikesCount -= 1
		}
	case types.LikeWeightZero:
		switch newWeight {
		case types.LikeWeightUp:
			post.LikesCount += 1
		case types.LikeWeightDown:
			post.DislikesCount -= 1
		}
	}
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

// GetPostsIterator returns an iterator over all likes
func (k Keeper) GetLikesIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.LikePrefix)
}

func getPostKeeperKey(p types.Post) []byte {
	return append(p.Owner.Bytes(), p.UUID[:]...)
}

func getLikeKeeperKey(l types.Like) []byte {
	return append(append(l.PostOwner.Bytes(), l.PostUUID[:]...), l.Owner.Bytes()...)
}

func getPostKeeperKeyFromLike(p types.Like) []byte {
	return append(p.PostOwner.Bytes(), p.PostUUID[:]...)
}
