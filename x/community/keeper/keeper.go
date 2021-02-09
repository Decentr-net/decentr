package keeper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"

	"github.com/Decentr-net/decentr/x/community/types"
	"github.com/Decentr-net/decentr/x/utils"
)

type Interval uint8

const (
	InvalidInterval Interval = iota
	DayInterval
	WeekInterval
	MonthInterval
)

var intervals = map[Interval]time.Duration{
	DayInterval:   time.Hour * 24,
	WeekInterval:  time.Hour * 24 * 7,
	MonthInterval: time.Hour * 24 * 31,
}

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	tokens   TokenKeeper

	paramSpace params.Subspace
}

// NewKeeper creates new instances of the community Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace, tokens TokenKeeper) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		tokens:     tokens,
		paramSpace: ps,
	}
}

// GetModerators returns the current moderators
func (k *Keeper) GetModerators(ctx sdk.Context) []string {
	var moderators []string
	k.paramSpace.GetIfExists(ctx, types.ParamModeratorsKey, &moderators)
	return moderators
}

// SetModerators sets the moderators
func (k *Keeper) SetModerators(ctx sdk.Context, moderators []string) {
	k.paramSpace.Set(ctx, types.ParamModeratorsKey, &moderators)
}

// CreatePost creates new post. Keeper's key is joined owner and uuid.
func (k Keeper) CreatePost(ctx sdk.Context, p types.Post) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	key := getPostKeeperKeyFromPost(p)
	store.Set(key, k.cdc.MustMarshalBinaryBare(p))

	createdAtIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexCreatedAtPrefix)
	indexKey := getCreateAtIndexKey(p)
	for _, p := range getCreatedAtIndexPrefixes(p.Category) {
		createdAtIndex.Set(append(p, indexKey...), key)
	}

	popularityIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexPopularPrefix)
	indexKey = getPopularityIndexKey(p)
	for _, p := range getPopularityIndexPrefixes(p.Category, p.CreatedAt, uint64(ctx.BlockTime().Unix())) {
		popularityIndex.Set(append(p, indexKey...), key)
	}
}

// DeletePost deletes the post from keeper.
func (k Keeper) DeletePost(ctx sdk.Context, owner sdk.AccAddress, id uuid.UUID) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	key := append(owner.Bytes(), id.Bytes()...)

	p := k.GetPostByKey(ctx, key)

	likesStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)
	it := sdk.KVStorePrefixIterator(likesStore, key)
	for ; it.Valid(); it.Next() {
		likesStore.Delete(it.Key())
	}
	it.Close()

	createdAtIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexCreatedAtPrefix)
	indexKey := getCreateAtIndexKey(p)
	for _, p := range getCreatedAtIndexPrefixes(p.Category) {
		createdAtIndex.Delete(append(p, indexKey...))
	}

	popularityIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexPopularPrefix)
	indexKey = getPopularityIndexKey(p)
	for _, p := range getPopularityIndexPrefixes(p.Category, 0, uint64(ctx.BlockTime().Unix())) {
		popularityIndex.Delete(append(p, indexKey...))
	}

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

	postKey := getPostKeeperKeyFromLike(newLike)

	postsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)
	oldPostBytes := postsStore.Get(postKey)

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
	k.updatePostPDV(ctx, &newPost, utils.GetHash(newLike))

	postsStore.Set(postKey, k.cdc.MustMarshalBinaryBare(newPost))
	likesStore.Set(getLikeKeeperKey(newLike), k.cdc.MustMarshalBinaryBare(newLike))

	// update popularity index
	popularityIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexPopularPrefix)
	indexKey := getPopularityIndexKey(oldPost)
	for _, p := range getPopularityIndexPrefixes(oldPost.Category, 0, uint64(ctx.BlockTime().Unix())) {
		popularityIndex.Delete(append(p, indexKey...))
	}

	indexKey = getPopularityIndexKey(newPost)
	for _, p := range getPopularityIndexPrefixes(newPost.Category, newPost.CreatedAt, uint64(ctx.BlockTime().Unix())) {
		popularityIndex.Set(append(p, indexKey...), postKey)
	}

	// update user likes index
	post := fmt.Sprintf("%s/%s", newLike.PostOwner, newLike.PostUUID)
	likes := make(map[string]types.LikeWeight)

	userLikesIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexUserLikesPrefix)
	if v := userLikesIndex.Get(newLike.Owner); v != nil {
		k.cdc.MustUnmarshalJSON(v, &likes)
	}

	if newLike.Weight == types.LikeWeightZero {
		delete(likes, post)
	} else {
		likes[post] = newLike.Weight
	}

	userLikesIndex.Set(newLike.Owner, sdk.MustSortJSON(k.cdc.MustMarshalJSON(likes)))
}

func (k Keeper) updatePostPDV(ctx sdk.Context, post *types.Post, description []byte) {
	oldPDV := post.PDV
	newPDV := sdk.NewInt(int64(post.LikesCount) - int64(post.DislikesCount))

	diff := newPDV.Sub(oldPDV)
	k.tokens.AddTokens(ctx, post.Owner, diff, description)

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
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix).Iterator(nil, nil)
}

func getPostKeeperKeyFromPost(p types.Post) []byte {
	return getPostKeeperKey(p.Owner, p.UUID)
}

func getPostKeeperKey(owner sdk.AccAddress, id uuid.UUID) []byte {
	return append(owner.Bytes(), id.Bytes()[:]...)
}

func getLikeKeeperKey(l types.Like) []byte {
	return append(append(l.PostOwner.Bytes(), l.PostUUID[:]...), l.Owner.Bytes()...)
}

func getPostKeeperKeyFromLike(p types.Like) []byte {
	return append(p.PostOwner.Bytes(), p.PostUUID[:]...)
}

func (k Keeper) SyncIndex(ctx sdk.Context) {
	popularityIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexPopularPrefix)

	flush := func(c types.Category, interval Interval, limit time.Duration) {
		l := uint64(ctx.BlockTime().Unix()) - uint64(limit/time.Second)

		index := prefix.NewStore(popularityIndex, []byte{byte(c), byte(interval)})

		it := sdk.KVStorePrefixIterator(index, nil)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			if k.GetPostByKey(ctx, it.Value()).CreatedAt < l {
				index.Delete(it.Key())
			}
		}
	}

	for i := types.UndefinedCategory; i <= types.SportsCategory; i++ {
		for b, d := range intervals {
			flush(i, b, d)
		}
	}
}

func (k Keeper) GetUserLikedPosts(ctx sdk.Context, owner sdk.AccAddress) map[string]types.LikeWeight {
	out := make(map[string]types.LikeWeight)

	userLikesIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexUserLikesPrefix)

	if userLikesIndex.Has(owner) {
		k.cdc.MustUnmarshalJSON(userLikesIndex.Get(owner), &out)
	}

	return out
}

func (k Keeper) GetRecentPosts(ctx sdk.Context, c types.Category, from []byte, limit uint32) []types.Post {
	return k.getPosts(ctx, append(types.IndexCreatedAtPrefix, byte(c)), from, limit)
}

func (k Keeper) GetPopularPosts(ctx sdk.Context, interval Interval, c types.Category, from []byte, limit uint32) []types.Post {
	return k.getPosts(ctx, append(types.IndexPopularPrefix, byte(c), byte(interval)), from, limit)
}

func (k Keeper) getPosts(ctx sdk.Context, p []byte, from []byte, limit uint32) []types.Post {
	index := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	it := sdk.KVStoreReversePrefixIterator(index, nil)
	defer it.Close()

	if from != nil {
		for ; it.Valid() && bytes.Compare(it.Key(), from) > -1; it.Next() {
		}
	}

	out := make([]types.Post, 0)
	for i := uint32(0); i < limit && it.Valid(); i++ {
		out = append(out, k.GetPostByKey(ctx, it.Value()))
		it.Next()
	}

	return out
}

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

func getCreatedAtIndexPrefixes(c types.Category) [][]byte {
	return [][]byte{
		{byte(types.UndefinedCategory)},
		{byte(c)},
	}
}

func getPopularityIndexPrefixes(c types.Category, createdAt uint64, timestamp uint64) [][]byte {
	catPrefixes := [][]byte{
		{byte(types.UndefinedCategory)},
		{byte(c)},
	}

	out := make([][]byte, 0, len(intervals)*len(catPrefixes))

	for _, p := range catPrefixes {
		for i, v := range intervals {
			if createdAt == 0 ||
				createdAt+uint64(v/time.Second) > timestamp {
				out = append(out, append(p, byte(i)))
			}
		}
	}

	return out
}

func getCreateAtIndexKey(p types.Post) []byte {
	return append(utils.Uint64ToBytes(p.CreatedAt), p.UUID.Bytes()...)
}

func getPopularityIndexKey(p types.Post) []byte {
	return append(utils.Uint64ToBytes(uint64(p.LikesCount)), p.UUID.Bytes()...)
}
