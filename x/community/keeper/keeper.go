package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"

	"github.com/Decentr-net/decentr/x/community/types"
)

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

	if err := k.index.AddPost(p); err != nil {
		ctx.Logger().Error("failed to add post to index",
			"err", err.Error(),
			"post", fmt.Sprintf("%s/%s", p.Owner, p.UUID.String()),
		)
	}

	store.Set(getPostKeeperKeyFromPost(p), k.cdc.MustMarshalBinaryBare(p))
}

// DeletePost deletes the post from keeper.
func (k Keeper) DeletePost(ctx sdk.Context, owner sdk.AccAddress, id uuid.UUID) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	key := append(owner.Bytes(), id.Bytes()...)

	if err := k.index.DeletePost(k.GetPostByKey(ctx, key)); err != nil {
		ctx.Logger().Error("failed to delete post from index",
			"err", err.Error(),
			"post", fmt.Sprintf("%s/%s", owner, id.String()),
		)
	}

	store.Delete(key)
}

// GetPostByKey returns entire post by keeper's key.
func (k Keeper) GetPostByKey(ctx sdk.Context, key []byte) types.Post {
	return k.getPostResolver(ctx)(key)
}

// GetPostsIterator returns an iterator over all posts
func (k Keeper) GetPostsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.PostPrefix)
}

func (k Keeper) ListUserPosts(ctx sdk.Context, owner sdk.AccAddress, from uuid.UUID, limit uint32) []types.Post {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	if from == uuid.Nil {
		// use max range
		from = uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	}

	it := store.ReverseIterator(getPostKeeperKey(owner, uuid.Nil), getPostKeeperKey(owner, from))
	defer it.Close()

	out := make([]types.Post, 0)

	for i := uint32(0); i < limit && it.Valid(); it.Next() {
		var post types.Post
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &post)
		out = append(out, post)

		i++
	}

	return out
}

func (k Keeper) SetLike(ctx sdk.Context, newLike types.Like) {
	if newLike.Owner.Equals(newLike.PostOwner) {
		ctx.Logger().Info("SetLike: owner tries to like own post",
			"postUUID", newLike.PostUUID,
			"postOwner", newLike.PostOwner)
		return
	}

	postsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)
	oldPostBytes := postsStore.Get(getPostKeeperKeyFromLike(newLike))

	if oldPostBytes == nil {
		ctx.Logger().Info("SetLike: post not found",
			"postUUID", newLike.PostUUID,
			"postOwner", newLike.PostOwner)
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
	k.updatePostPDV(ctx, &newPost)

	postsStore.Set(getPostKeeperKey(newPost.Owner, newPost.UUID), k.cdc.MustMarshalBinaryBare(newPost))
	likesStore.Set(getLikeKeeperKey(newLike), k.cdc.MustMarshalBinaryBare(newLike))
	_ = k.index.UpdateLikes(oldPost, newPost)
}

func (k Keeper) updatePostPDV(ctx sdk.Context, post *types.Post) {
	oldPDV := post.PDV
	newPDV := sdk.NewInt(int64(post.LikesCount) - int64(post.DislikesCount))

	diff := newPDV.Sub(oldPDV)
	k.tokens.AddTokens(ctx, post.Owner, post.CreatedAt, diff)

	post.PDV = newPDV
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
			post.DislikesCount += 1
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

// GetLikesIterator returns an iterator over all likes
func (k Keeper) GetLikesIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.LikePrefix)
}

func getPostKeeperKeyFromPost(p types.Post) []byte {
	return getPostKeeperKey(p.Owner, p.UUID)
}

func getPostKeeperKey(owner sdk.AccAddress, id uuid.UUID) []byte {
	return append(owner.Bytes(), id.Bytes()[:]...)
}

func (k Keeper) getPostResolver(ctx sdk.Context) func([]byte) types.Post {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	return func(key []byte) types.Post {
		if !store.Has(key) {
			return types.Post{}
		}

		bz := store.Get(key)

		var post types.Post
		k.cdc.MustUnmarshalBinaryBare(bz, &post)
		return post
	}
}

func getLikeKeeperKey(l types.Like) []byte {
	return append(append(l.PostOwner.Bytes(), l.PostUUID[:]...), l.Owner.Bytes()...)
}

func getPostKeeperKeyFromLike(p types.Like) []byte {
	return append(p.PostOwner.Bytes(), p.PostUUID[:]...)
}

func (k Keeper) SyncIndex(ctx sdk.Context) {
	k.index.RemoveUnnecessaryPosts(ctx, uint64(ctx.BlockTime().Unix()), k.getPostResolver(ctx))
}
