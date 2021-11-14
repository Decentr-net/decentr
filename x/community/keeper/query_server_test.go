package keeper

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/community/types"
)

func TestQueryServer_Moderators(t *testing.T) {
	set, ctx := setupKeeper(t)
	p := types.Params{
		Moderators: []sdk.AccAddress{NewAccAddress(), NewAccAddress()},
		FixedGas:   types.FixedGasParams{},
	}
	set.keeper.SetParams(ctx, p)
	s := NewQueryServer(set.keeper)

	resp, err := s.Moderators(sdk.WrapSDKContext(ctx), nil)
	require.NoError(t, err)
	require.Equal(t, p.Moderators, resp.Moderators)
}

func TestQueryServer_GetPost(t *testing.T) {
	set, ctx := setupKeeper(t)
	p := types.Post{
		Owner:        NewAccAddress(),
		Uuid:         uuid.Must(uuid.NewV1()).String(),
		Title:        "title",
		PreviewImage: "http://decentr.xyz/preview.png",
		Category:     types.Category_CATEGORY_CRYPTO_AND_BLOCKCHAIN,
		Text:         "fifteen symbols should be here",
	}
	set.keeper.SetPost(ctx, p)
	s := NewQueryServer(set.keeper)

	t.Run("success", func(t *testing.T) {
		resp, err := s.GetPost(sdk.WrapSDKContext(ctx), &types.GetPostRequest{
			PostOwner: p.Owner,
			PostUuid:  p.Uuid,
		})
		require.NoError(t, err)
		require.Equal(t, p, resp.Post)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := s.GetPost(sdk.WrapSDKContext(ctx), &types.GetPostRequest{
			PostOwner: p.Owner,
			PostUuid:  "uuid",
		})
		require.Error(t, err)
		require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrInvalidRequest))
	})
}

func TestQueryServer_ListUserPosts(t *testing.T) {
	set, ctx := setupKeeper(t)

	owner := NewAccAddress()
	posts := make([]types.Post, 10)
	for i := range posts {
		p := types.Post{
			Owner:        owner,
			Uuid:         uuid.Must(uuid.NewV1()).String(),
			Title:        "title",
			PreviewImage: "http://decentr.xyz/preview.png",
			Category:     types.Category_CATEGORY_CRYPTO_AND_BLOCKCHAIN,
			Text:         "fifteen symbols should be here",
		}
		set.keeper.SetPost(ctx, p)
		posts[i] = p
	}

	s := NewQueryServer(set.keeper)

	t.Run("ok default", func(t *testing.T) {
		resp, err := s.ListUserPosts(sdk.WrapSDKContext(ctx), &types.ListUserPostsRequest{
			Owner:      owner,
			Pagination: query.PageRequest{},
		})
		require.NoError(t, err)
		require.ElementsMatch(t, posts, resp.Posts)
	})

	t.Run("limited", func(t *testing.T) {
		resp, err := s.ListUserPosts(sdk.WrapSDKContext(ctx), &types.ListUserPostsRequest{
			Owner: owner,
			Pagination: query.PageRequest{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Posts, 1)
	})
}

func TestQueryServer_ListFollowed(t *testing.T) {
	set, ctx := setupKeeper(t)

	owner := NewAccAddress()
	followed := make([]sdk.AccAddress, 10)
	for i := range followed {
		followed[i] = NewAccAddress()
		set.keeper.Follow(ctx, owner, followed[i])
	}

	s := NewQueryServer(set.keeper)

	t.Run("ok default", func(t *testing.T) {
		resp, err := s.ListFollowed(sdk.WrapSDKContext(ctx), &types.ListFollowedRequest{
			Owner:      owner,
			Pagination: query.PageRequest{},
		})
		require.NoError(t, err)
		require.ElementsMatch(t, followed, resp.Followed)
	})

	t.Run("limited", func(t *testing.T) {
		resp, err := s.ListFollowed(sdk.WrapSDKContext(ctx), &types.ListFollowedRequest{
			Owner: owner,
			Pagination: query.PageRequest{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Followed, 1)
	})
}
