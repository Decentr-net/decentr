package keeper

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/community/types"
)

func TestMsgServer_CreatePost(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.tokenKeeper)

	owner := NewAccAddress()
	id := uuid.Must(uuid.NewV1())

	_, err := s.CreatePost(sdk.WrapSDKContext(ctx), &types.MsgCreatePost{
		Post: types.Post{
			Owner:        owner.String(),
			Uuid:         id.String(),
			Title:        "gr8 title",
			PreviewImage: "",
			Category:     1,
			Text:         "Fifteen symbols have to be here",
		},
	})
	require.NoError(t, err)

	_, err = s.CreatePost(sdk.WrapSDKContext(ctx), &types.MsgCreatePost{
		Post: types.Post{
			Owner:        owner.String(),
			Uuid:         id.String(),
			Title:        "gr8 title",
			PreviewImage: "",
			Category:     1,
			Text:         "Fifteen symbols have to be here",
		},
	})
	require.Error(t, err)
	require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrConflict))
}

func TestMsgServer_DeletePost(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.tokenKeeper)

	owner := NewAccAddress()
	id := uuid.Must(uuid.NewV1())

	createPost := func() {
		_, err := s.CreatePost(sdk.WrapSDKContext(ctx), &types.MsgCreatePost{
			Post: types.Post{
				Owner:        owner.String(),
				Uuid:         id.String(),
				Title:        "gr8 title",
				PreviewImage: "",
				Category:     1,
				Text:         "Fifteen symbols have to be here",
			},
		})
		require.NoError(t, err)
		require.NotEqual(t, types.Post{}, set.keeper.GetPost(ctx, postKey(owner, id)))
	}

	t.Run("forbidden", func(t *testing.T) {
		_, err := s.DeletePost(sdk.WrapSDKContext(ctx), &types.MsgDeletePost{
			Owner:     NewAccAddress().String(),
			PostOwner: owner.String(),
			PostUuid:  id.String(),
		})
		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrUnauthorized))
	})

	t.Run("not found", func(t *testing.T) {
		_, err := s.DeletePost(sdk.WrapSDKContext(ctx), &types.MsgDeletePost{
			Owner:     owner.String(),
			PostOwner: owner.String(),
			PostUuid:  id.String(),
		})
		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrNotFound))
	})

	t.Run("deleted by owner", func(t *testing.T) {
		createPost()
		_, err := s.DeletePost(sdk.WrapSDKContext(ctx), &types.MsgDeletePost{
			Owner:     owner.String(),
			PostOwner: owner.String(),
			PostUuid:  id.String(),
		})
		require.NoError(t, err)
		require.Equal(t, types.Post{}, set.keeper.GetPost(ctx, postKey(owner, id)))
	})

	t.Run("delete by moderator", func(t *testing.T) {
		createPost()
		moderator := NewAccAddress()
		set.keeper.SetParams(ctx, types.Params{
			Moderators: []string{moderator.String()},
			FixedGas:   types.DefaultFixedGasParams(),
		})

		_, err := s.DeletePost(sdk.WrapSDKContext(ctx), &types.MsgDeletePost{
			Owner:     moderator.String(),
			PostOwner: owner.String(),
			PostUuid:  id.String(),
		})
		require.NoError(t, err)
		require.Equal(t, types.Post{}, set.keeper.GetPost(ctx, postKey(owner, id)))
	})
}

func TestMsgServer_SetLike(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.tokenKeeper)

	owner := NewAccAddress()
	id := uuid.Must(uuid.NewV1())
	createPost := func() {
		_, err := s.CreatePost(sdk.WrapSDKContext(ctx), &types.MsgCreatePost{
			Post: types.Post{
				Owner:        owner.String(),
				Uuid:         id.String(),
				Title:        "gr8 title",
				PreviewImage: "",
				Category:     1,
				Text:         "Fifteen symbols have to be here",
			},
		})
		require.NoError(t, err)
		require.NotEqual(t, types.Post{}, set.keeper.GetPost(ctx, postKey(owner, id)))
	}

	t.Run("not found", func(t *testing.T) {
		_, err := s.SetLike(sdk.WrapSDKContext(ctx), &types.MsgSetLike{
			Like: types.Like{
				Owner:     NewAccAddress().String(),
				PostOwner: NewAccAddress().String(),
				PostUuid:  id.String(),
			},
		})
		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrNotFound))
	})

	t.Run("self like", func(t *testing.T) {
		_, err := s.SetLike(sdk.WrapSDKContext(ctx), &types.MsgSetLike{
			Like: types.Like{
				Owner:     owner.String(),
				PostOwner: owner.String(),
				PostUuid:  id.String(),
			},
		})
		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrInvalidRequest))
	})

	t.Run("like", func(t *testing.T) {
		likeOwner := NewAccAddress()

		createPost()
		_, err := s.SetLike(sdk.WrapSDKContext(ctx), &types.MsgSetLike{
			Like: types.Like{
				Owner:     likeOwner.String(),
				PostOwner: owner.String(),
				PostUuid:  id.String(),
				Weight:    types.LikeWeight_LIKE_WEIGHT_UP,
			},
		})
		require.NoError(t, err)

		_, err = s.SetLike(sdk.WrapSDKContext(ctx), &types.MsgSetLike{
			Like: types.Like{
				Owner:     likeOwner.String(),
				PostOwner: owner.String(),
				PostUuid:  id.String(),
				Weight:    types.LikeWeight_LIKE_WEIGHT_DOWN,
			},
		})
		require.NoError(t, err)
	})
}

func TestMsgServer_Follow(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.tokenKeeper)

	t.Run("self like", func(t *testing.T) {
		address := NewAccAddress()
		_, err := s.Follow(sdk.WrapSDKContext(ctx), &types.MsgFollow{
			Owner: address.String(),
			Whom:  address.String(),
		})

		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrInvalidRequest))
	})

	t.Run("already follows", func(t *testing.T) {
		who, whom := NewAccAddress(), NewAccAddress()
		set.keeper.Follow(ctx, who, whom)
		_, err := s.Follow(sdk.WrapSDKContext(ctx), &types.MsgFollow{
			Owner: who.String(),
			Whom:  whom.String(),
		})

		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrConflict))
	})

	t.Run("correct", func(t *testing.T) {
		_, err := s.Follow(sdk.WrapSDKContext(ctx), &types.MsgFollow{
			Owner: NewAccAddress().String(),
			Whom:  NewAccAddress().String(),
		})

		require.NoError(t, err)
	})
}

func TestMsgServer_Unfollow(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.tokenKeeper)

	who, whom := NewAccAddress(), NewAccAddress()

	t.Run("not found", func(t *testing.T) {
		_, err := s.Unfollow(sdk.WrapSDKContext(ctx), &types.MsgUnfollow{
			Owner: who.String(),
			Whom:  whom.String(),
		})

		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrNotFound))
	})

	t.Run("correct", func(t *testing.T) {
		set.keeper.Follow(ctx, who, whom)
		_, err := s.Unfollow(sdk.WrapSDKContext(ctx), &types.MsgUnfollow{
			Owner: who.String(),
			Whom:  whom.String(),
		})

		require.NoError(t, err)
	})
}
