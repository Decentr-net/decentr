package keeper

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	tokenkeeper "github.com/Decentr-net/decentr/x/token/keeper"
	tokentypes "github.com/Decentr-net/decentr/x/token/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/community/types"
)

func init() {
	config.SetAddressPrefixes()
}

type keeperSet struct {
	keeper      Keeper
	tokenKeeper types.TokenKeeper
}

func setupKeeper(t testing.TB) (keeperSet, sdk.Context) {
	keys := sdk.NewKVStoreKeys(types.StoreKey, tokentypes.StoreKey, paramstypes.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	ctx, err := GetContext(keys, tkeys)
	require.NoError(t, err)

	set := keeperSet{}

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsKeeper := paramskeeper.NewKeeper(
		cdc,
		codec.NewLegacyAmino(),
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	tokenKeeper := tokenkeeper.NewKeeper(
		cdc,
		keys[tokentypes.StoreKey],
	)
	set.tokenKeeper = tokenKeeper

	set.keeper = *NewKeeper(
		cdc,
		keys[types.StoreKey],
		paramsKeeper.Subspace(types.StoreKey),
	)
	set.keeper.SetParams(ctx, types.DefaultParams())

	return set, ctx
}

func TestKeeper_SetParams(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	require.NotPanics(t, func() { k.SetParams(ctx, types.DefaultParams()) })
}

func TestKeeper_GetParams(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	require.Equal(t, types.DefaultParams(), k.GetParams(ctx))

	p := types.Params{
		Moderators: []string{NewAccAddress().String()},
		FixedGas: types.FixedGasParams{
			CreatePost: 1,
			DeletePost: 2,
			SetLike:    3,
			Follow:     4,
			Unfollow:   5,
		},
	}

	k.SetParams(ctx, p)
	require.Equal(t, p, k.GetParams(ctx))
}

func TestKeeper_IsModerator(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	a1, a2, a3 := NewAccAddress(), NewAccAddress(), NewAccAddress()

	p := types.DefaultParams()
	p.Moderators = []string{a1.String(), a2.String()}
	k.SetParams(ctx, p)

	require.True(t, k.IsModerator(ctx, a1))
	require.True(t, k.IsModerator(ctx, a2))
	require.False(t, k.IsModerator(ctx, a3))
}

func TestKeeper_parseLikeKey(t *testing.T) {
	addr1, id, addr2 := NewAccAddress(), uuid.Must(uuid.NewV1()), NewAccAddress()

	actAddr1, actID, actAddr2 := parseLikeKey(likeKey(postKey(addr1, id), addr2))

	require.Equal(t, addr1, actAddr1)
	require.Equal(t, id, actID)
	require.Equal(t, addr2, actAddr2)
}

func TestKeeper_SetPost(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	owner := NewAccAddress()
	id := uuid.Must(uuid.NewV1())

	p := types.Post{
		Owner:        owner.String(),
		Uuid:         id.String(),
		Title:        "title",
		PreviewImage: "image",
		Category:     1,
		Text:         "text",
	}

	k.SetPost(ctx, p)
	require.Equal(t, p, k.GetPost(ctx, postKey(owner, id)))
}

func TestKeeper_GetPost(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	require.Equal(t, types.Post{}, k.GetPost(ctx, postKey(NewAccAddress(), uuid.Must(uuid.NewV1()))))
}

func TestKeeper_DeletePost(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	owner := NewAccAddress()
	id := uuid.Must(uuid.NewV1())

	p := types.Post{
		Owner:        owner.String(),
		Uuid:         id.String(),
		Title:        "title",
		PreviewImage: "image",
		Category:     1,
		Text:         "text",
	}

	k.SetPost(ctx, p)

	key := postKey(owner, id)
	k.DeletePost(ctx, key)
	require.Equal(t, types.Post{}, k.GetPost(ctx, key))
}

func TestKeeper_IteratePosts(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	exp := make([]types.Post, 10)
	for i := 0; i < 10; i++ {
		exp[i] = types.Post{
			Owner:        NewAccAddress().String(),
			Uuid:         uuid.Must(uuid.NewV1()).String(),
			Title:        "title",
			PreviewImage: "image",
			Category:     1,
			Text:         "text",
		}
		k.SetPost(ctx, exp[i])
	}

	act := make([]types.Post, 0, 10)
	k.IteratePosts(ctx, func(p types.Post) (stop bool) {
		act = append(act, p)
		return
	})

	require.ElementsMatch(t, exp, act)
}

func TestKeeper_ListUserPosts(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	owner := NewAccAddress()

	exp := make([]types.Post, 10)
	for i := 0; i < 10; i++ {
		exp[i] = types.Post{
			Owner: owner.String(),
			Uuid:  uuid.Must(uuid.NewV1()).String(),
			Title: fmt.Sprintf("title #%d", i),
		}
		k.SetPost(ctx, exp[i])
	}

	t.Run("fetch all", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListUserPosts(ctx, owner, query.PageRequest{
			Offset:     0,
			Limit:      10,
			CountTotal: true,
		})
		require.Equal(t, exp, act)
		require.Equal(t, uuid.FromStringOrNil(exp[9].Uuid).Bytes(), next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListUserPosts(ctx, owner, query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: true,
		})
		require.Equal(t, exp[0], act[0])
		require.Equal(t, uuid.FromStringOrNil(exp[0].Uuid).Bytes(), next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit reverse", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListUserPosts(ctx, owner, query.PageRequest{
			Offset:     0,
			Limit:      1,
			Reverse:    true,
			CountTotal: true,
		})
		require.Equal(t, exp[9], act[0])
		require.Equal(t, uuid.FromStringOrNil(exp[9].Uuid).Bytes(), next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit offset", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListUserPosts(ctx, owner, query.PageRequest{
			Offset:     1,
			Limit:      1,
			CountTotal: true,
		})
		require.Equal(t, exp[1], act[0])
		require.Equal(t, uuid.FromStringOrNil(exp[1].Uuid).Bytes(), next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit key", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListUserPosts(ctx, owner, query.PageRequest{
			Key:        uuid.FromStringOrNil(exp[1].Uuid).Bytes(),
			Limit:      1,
			CountTotal: true,
		})
		require.Equal(t, exp[2], act[0])
		require.Equal(t, uuid.FromStringOrNil(exp[2].Uuid).Bytes(), next)
		require.Len(t, exp, int(total))
	})
}

func TestKeeper_SetLike(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	owner, postOwner := NewAccAddress(), NewAccAddress()
	id := uuid.Must(uuid.NewV1())

	l := types.Like{
		Owner:     owner.String(),
		PostOwner: postOwner.String(),
		PostUuid:  id.String(),
		Weight:    types.LikeWeight_LIKE_WEIGHT_UP,
	}

	k.SetLike(ctx, l)
	require.Equal(t, l, k.GetLike(ctx, likeKey(postKey(postOwner, id), owner)))

	l.Weight = types.LikeWeight_LIKE_WEIGHT_ZERO
	k.SetLike(ctx, l)
	require.Equal(t, types.Like{}, k.GetLike(ctx, likeKey(postKey(postOwner, id), owner)))
}

func TestKeeper_GetLike(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	require.Equal(t, types.Like{}, k.GetLike(
		ctx,
		likeKey(postKey(NewAccAddress(), uuid.Must(uuid.NewV1())), NewAccAddress())),
	)
}

func TestKeeper_IterateLikes(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	exp := make([]types.Like, 10)
	for i := 0; i < 10; i++ {
		exp[i] = types.Like{
			Owner:     NewAccAddress().String(),
			PostOwner: NewAccAddress().String(),
			PostUuid:  uuid.Must(uuid.NewV1()).String(),
			Weight:    []types.LikeWeight{-1, 1}[i%2],
		}
		k.SetLike(ctx, exp[i])
	}

	act := make([]types.Like, 0, 10)
	k.IterateLikes(ctx, func(p types.Like) (stop bool) {
		act = append(act, p)
		return
	})

	require.ElementsMatch(t, exp, act)
}

func TestKeeper_Follow(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	addr1, addr2 := NewAccAddress(), NewAccAddress()
	k.Follow(ctx, addr1, addr2)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowingPrefix)
	require.True(t, store.Has(append(addr1, addr2...)))
}

func TestKeeper_Unfollow(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	addr1, addr2 := NewAccAddress(), NewAccAddress()
	k.Follow(ctx, addr1, addr2)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FollowingPrefix)

	k.Unfollow(ctx, addr2, addr1)
	require.True(t, store.Has(append(addr1, addr2...)))
	k.Unfollow(ctx, addr1, addr2)
	require.False(t, store.Has(append(addr1, addr2...)))
}

func TestKeeper_IterateFollowings(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	type pair struct {
		who  sdk.AccAddress
		whom sdk.AccAddress
	}
	exp := make([]pair, 0)
	for i := 0; i < 10; i++ {
		who := NewAccAddress()
		for j := 0; j < 10; j++ {
			whom := NewAccAddress()
			exp = append(exp, pair{
				who:  who,
				whom: whom,
			})
			k.Follow(ctx, who, whom)
		}
	}

	act := make([]pair, 0)
	k.IterateFollowings(ctx, func(who, whom sdk.AccAddress) (stop bool) {
		act = append(act, pair{
			who:  who,
			whom: whom,
		})
		return
	})

	require.ElementsMatch(t, exp, act)
}

func TestKeeper_IsFollowed(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	who, whom := NewAccAddress(), NewAccAddress()
	require.False(t, k.IsFollowed(ctx, who, whom))

	k.Follow(ctx, who, whom)
	require.True(t, k.IsFollowed(ctx, who, whom))
}

func TestKeeper_ListFollowed(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	owner := NewAccAddress()

	exp := make([]sdk.AccAddress, 10)
	for i := 0; i < 10; i++ {
		exp[i] = NewAccAddress()
		k.Follow(ctx, owner, exp[i])
	}

	ordered, _, _ := k.ListFollowed(ctx, owner, query.PageRequest{
		Limit: 10,
	})
	require.ElementsMatch(t, exp, ordered)
	exp = ordered

	t.Run("fetch all", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListFollowed(ctx, owner, query.PageRequest{
			Offset:     0,
			Limit:      10,
			CountTotal: true,
		})
		require.Equal(t, exp, act)
		require.EqualValues(t, exp[9], next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListFollowed(ctx, owner, query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: true,
		})
		require.Equal(t, exp[0], act[0])
		require.EqualValues(t, exp[0], next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit reverse", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListFollowed(ctx, owner, query.PageRequest{
			Offset:     0,
			Limit:      1,
			Reverse:    true,
			CountTotal: true,
		})
		require.Equal(t, exp[9], act[0])
		require.EqualValues(t, exp[9], next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit offset", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListFollowed(ctx, owner, query.PageRequest{
			Offset:     1,
			Limit:      1,
			CountTotal: true,
		})
		require.Equal(t, exp[1], act[0])
		require.EqualValues(t, exp[1], next)
		require.Len(t, exp, int(total))
	})
	t.Run("limit key", func(t *testing.T) {
		t.Parallel()

		act, next, total := k.ListFollowed(ctx, owner, query.PageRequest{
			Key:        exp[1],
			Limit:      1,
			CountTotal: true,
		})
		require.Equal(t, exp[2], act[0])
		require.EqualValues(t, exp[2], next)
		require.Len(t, exp, int(total))
	})
}

func TestKeeper_ResetAccount(t *testing.T) {
	set, ctx := setupKeeper(t)
	k := set.keeper

	owner := NewAccAddress()
	addr := NewAccAddress()
	id := uuid.Must(uuid.NewV1())

	k.SetPost(ctx, types.Post{
		Owner:        owner.String(),
		Uuid:         id.String(),
		Title:        "title",
		PreviewImage: "image",
		Category:     1,
		Text:         "text",
	})

	k.SetLike(ctx, types.Like{
		Owner:     addr.String(),
		PostOwner: owner.String(),
		PostUuid:  id.String(),
		Weight:    1,
	})

	k.Follow(ctx, addr, owner)
	k.Follow(ctx, owner, addr)

	k.ResetAccount(ctx, owner)

	// Should be nothing
	k.IteratePosts(ctx, func(_ types.Post) (stop bool) {
		t.Fatal("account has posts")
		return
	})
	k.IterateLikes(ctx, func(l types.Like) (stop bool) {
		t.Fatal("account's post is still liked")
		return
	})
	k.IterateFollowings(ctx, func(who, whom sdk.AccAddress) (stop bool) {
		t.Fatal("account follows or is followed")
		return
	})
}
