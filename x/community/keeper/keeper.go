package keeper

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/gofrs/uuid"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Decentr-net/decentr/x/community/types"
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

func postKey(owner sdk.AccAddress, id uuid.UUID) []byte {
	return append(owner, id.Bytes()...)
}

func likeKey(postKey []byte, owner sdk.AccAddress) []byte {
	return append(postKey, owner...)
}

func parseLikeKey(k []byte) (postOwner sdk.AccAddress, postUUID uuid.UUID, owner sdk.AccAddress) {
	// calculate addr len
	addrLen := (len(k) - uuid.Size) / 2

	postOwner, k = k[:addrLen], k[addrLen:]
	if err := sdk.VerifyAddressFormat(postOwner); err != nil {
		panic(fmt.Errorf("invalid like key: invalid post_owner: %w", err))
	}

	if err := postUUID.UnmarshalBinary(k[:uuid.Size]); err != nil {
		panic("invalid like key")
	}

	owner = k[uuid.Size:]
	if err := sdk.VerifyAddressFormat(postOwner); err != nil {
		panic(fmt.Errorf("invalid like key: invalid owner: %w", err))
	}

	return
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

func (k Keeper) IsModerator(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, v := range k.GetParams(ctx).Moderators {
		addr, _ := sdk.AccAddressFromBech32(v)
		if address.Equals(addr) {
			return true
		}
	}

	return false
}

func (k Keeper) SetPost(ctx sdk.Context, p types.Post) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	id, err := uuid.FromString(p.Uuid)
	if err != nil {
		panic(fmt.Sprintf("invalid uuid in post: %+v", p))
	}

	bz, err := p.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal post %s/%s: %w", p.Owner, p.Uuid, err))
	}

	owner, _ := sdk.AccAddressFromBech32(p.Owner)
	store.Set(postKey(owner, id), bz)
}

func (k Keeper) DeletePost(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	if store.Has(key) {
		store.Delete(key)
	}

	k.dropPostLikes(ctx, key)
}

func (k Keeper) GetPost(ctx sdk.Context, key []byte) types.Post {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	if !k.HasPost(ctx, key) {
		return types.Post{}
	}

	var p types.Post
	if err := p.Unmarshal(store.Get(key)); err != nil {
		panic(fmt.Errorf("failed to get post: %w", err))
	}

	return p
}

func (k Keeper) HasPost(ctx sdk.Context, key []byte) bool {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix).Has(key)
}

func (k Keeper) ListUserPosts(
	ctx sdk.Context,
	owner sdk.AccAddress,
	p query.PageRequest,
) (posts []types.Post, nextKey []byte, total uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.PostPrefix, owner...))

	var it sdk.Iterator
	if !p.Reverse {
		it = store.Iterator(nil, nil)
	} else {
		it = store.ReverseIterator(nil, nil)
	}
	defer it.Close()

	for it.Valid() {
		key := it.Key()
		total++
		it.Next()

		// skip before offset is reached
		if total <= p.Offset {
			continue
		}

		// skip all keys before cursor
		if p.Key != nil {
			if bytes.Equal(key, p.Key) {
				p.Key = nil
			}
			continue
		}

		if uint64(len(posts)) >= p.Limit {
			if !p.CountTotal {
				return posts, nextKey, 0
			}

			continue
		}

		posts = append(posts, k.GetPost(ctx, append(owner, key...)))
		nextKey = key
	}

	return posts, nextKey, total
}

func (k Keeper) GetLike(ctx sdk.Context, key []byte) types.Like {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)

	if !store.Has(key) {
		return types.Like{}
	}

	postOwner, postUUID, owner := parseLikeKey(key)

	return types.Like{
		Owner:     owner.String(),
		PostOwner: postOwner.String(),
		PostUuid:  postUUID.String(),
		Weight:    types.LikeWeight(int8(store.Get(key)[0])),
	}
}

func (k Keeper) SetLike(ctx sdk.Context, l types.Like) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)

	postUUID, err := uuid.FromString(l.PostUuid)
	if err != nil {
		panic(fmt.Sprintf("invalid uuid in like: %+v", l))
	}

	owner, _ := sdk.AccAddressFromBech32(l.Owner)
	postOwner, _ := sdk.AccAddressFromBech32(l.PostOwner)

	key := likeKey(postKey(postOwner, postUUID), owner)

	if l.Weight == types.LikeWeight_LIKE_WEIGHT_ZERO {
		store.Delete(key)
	} else {
		store.Set(key, []byte{byte(l.Weight)})
	}
}

func (k Keeper) dropPostLikes(ctx sdk.Context, postKey []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.LikePrefix, postKey...))

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		store.Delete(it.Key())
		it.Next()
	}
}

func (k Keeper) Follow(ctx sdk.Context, who, whom sdk.AccAddress) {
	if who.Equals(whom) {
		return
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.FollowingPrefix, who...))

	store.Set(whom, []byte{})
}

func (k Keeper) Unfollow(ctx sdk.Context, who, whom sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.FollowingPrefix, who...))

	store.Delete(whom)
}

func (k Keeper) IsFollowed(ctx sdk.Context, who, whom sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.FollowingPrefix, who...))

	return store.Has(whom)
}

func (k Keeper) ListFollowed(
	ctx sdk.Context,
	owner sdk.AccAddress,
	p query.PageRequest,
) (followed []sdk.AccAddress, nextKey []byte, total uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.FollowingPrefix, owner...))

	var it sdk.Iterator
	if p.Reverse {
		it = store.Iterator(nil, nil)
	} else {
		it = store.ReverseIterator(nil, nil)
	}
	defer it.Close()

	for it.Valid() {
		key := it.Key()
		total++
		it.Next()

		// skip before offset is reached
		if total <= p.Offset {
			continue
		}

		// skip all keys before cursor
		if p.Key != nil {
			if bytes.Equal(key, p.Key) {
				p.Key = nil
			}
			continue
		}

		if uint64(len(followed)) >= p.Limit {
			if !p.CountTotal {
				return followed, nextKey, 0
			}

			continue
		}

		followed = append(followed, key)
		nextKey = key
	}

	return followed, nextKey, total
}

func (k Keeper) ResetAccount(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.PostPrefix, owner...))

	it := store.Iterator(nil, nil)
	for it.Valid() {
		k.DeletePost(ctx, append(owner, it.Key()...))
		it.Next()
	}
	it.Close()

	k.IterateFollowings(ctx, func(who, whom sdk.AccAddress) (stop bool) {
		if owner.Equals(who) || owner.Equals(whom) {
			k.Unfollow(ctx, who, whom)
		}
		return false
	})
}

func (k Keeper) IteratePosts(ctx sdk.Context, handle func(p types.Post) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PostPrefix)

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		if handle(k.GetPost(ctx, it.Key())) {
			return
		}

		it.Next()
	}
}

func (k Keeper) IterateLikes(ctx sdk.Context, handle func(l types.Like) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LikePrefix)

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		if handle(k.GetLike(ctx, it.Key())) {
			return
		}

		it.Next()
	}
}

func (k Keeper) IterateFollowings(ctx sdk.Context, handle func(who, whom sdk.AccAddress) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowingPrefix)

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		addrLen := len(it.Key()) / 2
		who := it.Key()[:addrLen]
		whom := it.Key()[addrLen:]

		if handle(who, whom) {
			return
		}

		it.Next()
	}
}
