package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/gofrs/uuid"

	"github.com/Decentr-net/decentr/x/community/types"
)

var _ types.QueryServer = queryServer{}

const defaultLimit = 50

type queryServer struct {
	keeper Keeper
}

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

func (s queryServer) GetPost(
	goCtx context.Context,
	r *types.GetPostRequest,
) (*types.GetPostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id, err := uuid.FromString(r.PostUuid)
	if err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest, "invalid post_uuid",
		)
	}

	postOwner, err := sdk.AccAddressFromBech32(r.PostOwner)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid post_owner address: %s", err)
	}

	return &types.GetPostResponse{
		Post: s.keeper.GetPost(ctx, postKey(postOwner, id)),
	}, nil
}

func (s queryServer) ListUserPosts(
	goCtx context.Context,
	r *types.ListUserPostsRequest,
) (*types.ListUserPostsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if r.Pagination.Key != nil && r.Pagination.Offset != 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid pagination: only one of offset or key can be set")
	}

	if r.Pagination.Limit == 0 {
		r.Pagination.Limit = defaultLimit
	}

	owner, err := sdk.AccAddressFromBech32(r.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	posts, next, total := s.keeper.ListUserPosts(ctx, owner, r.Pagination)

	return &types.ListUserPostsResponse{
		Posts: posts,
		Pagination: query.PageResponse{
			NextKey: next,
			Total:   total,
		},
	}, nil
}

func (s queryServer) Moderators(goCtx context.Context, _ *types.ModeratorsRequest) (*types.ModeratorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.ModeratorsResponse{
		Moderators: s.keeper.GetParams(ctx).Moderators,
	}, nil
}

func (s queryServer) ListFollowed(
	goCtx context.Context,
	r *types.ListFollowedRequest,
) (*types.ListFollowedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if r.Pagination.Key != nil && r.Pagination.Offset != 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid pagination: only one of offset or key can be set")
	}

	if r.Pagination.Limit == 0 {
		r.Pagination.Limit = defaultLimit
	}

	owner, err := sdk.AccAddressFromBech32(r.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	followed, next, total := s.keeper.ListFollowed(ctx, owner, r.Pagination)

	return &types.ListFollowedResponse{
		Followed: func() []string {
			out := make([]string, len(followed))
			for i, v := range followed {
				out[i] = v.String()
			}
			return out
		}(),
		Pagination: query.PageResponse{
			NextKey: next,
			Total:   total,
		},
	}, nil
}
