package keeper

import (
	"context"

	"github.com/gofrs/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/community/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper Keeper

	tokenKeeper types.TokenKeeper
}

// NewMsgServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServer(
	keeper Keeper,
	tokenKeeper types.TokenKeeper,
) types.MsgServer {
	return &msgServer{
		keeper:      keeper,
		tokenKeeper: tokenKeeper,
	}
}

func (s msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id, err := uuid.FromString(msg.Post.Uuid)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post_uuid: %s", err.Error())
	}

	if s.keeper.HasPost(ctx, postKey(msg.Post.Owner, id)) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrConflict, "post %s already exists", msg.Post.Address())
	}

	s.keeper.SetPost(ctx, msg.Post)

	return &types.MsgCreatePostResponse{}, nil
}

func (s msgServer) DeletePost(goCtx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !msg.Owner.Equals(msg.PostOwner) {
		if !s.keeper.IsModerator(ctx, msg.Owner) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a moderator or post owner", msg.Owner)
		}
	}

	id, err := uuid.FromString(msg.PostUuid)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post_uuid: %s", err.Error())
	}

	key := postKey(msg.PostOwner, id)

	if !s.keeper.HasPost(ctx, postKey(msg.PostOwner, id)) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "post %s/%s not found", msg.PostOwner, msg.PostUuid)
	}

	s.keeper.DeletePost(ctx, key)

	return &types.MsgDeletePostResponse{}, nil
}

func (s msgServer) SetLike(goCtx context.Context, msg *types.MsgSetLike) (*types.MsgSetLikeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	postUUID, err := uuid.FromString(msg.Like.PostUuid)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post_uuid: %s", err.Error())
	}

	if msg.Like.Owner.Equals(msg.Like.PostOwner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "self-like is disabled")
	}

	if !s.keeper.HasPost(ctx, postKey(msg.Like.PostOwner, postUUID)) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "post %s/%s not found", msg.Like.PostOwner, msg.Like.PostUuid)
	}

	diff := msg.Like.Weight - s.keeper.GetLike(ctx, likeKey(postKey(msg.Like.PostOwner, postUUID), msg.Like.Owner)).Weight
	s.tokenKeeper.IncTokens(ctx, msg.Like.PostOwner, sdk.NewDec(int64(diff)))

	s.keeper.SetLike(ctx, msg.Like)

	return &types.MsgSetLikeResponse{}, nil
}

func (s msgServer) Follow(goCtx context.Context, msg *types.MsgFollow) (*types.MsgFollowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Owner.Equals(msg.Whom) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "self-follow is disabled")
	}

	if s.keeper.IsFollowed(ctx, msg.Owner, msg.Whom) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrConflict, "%s already follows %s", msg.Owner, msg.Whom)
	}

	s.keeper.Follow(ctx, msg.Owner, msg.Whom)

	return &types.MsgFollowResponse{}, nil
}

func (s msgServer) Unfollow(goCtx context.Context, msg *types.MsgUnfollow) (*types.MsgUnfollowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !s.keeper.IsFollowed(ctx, msg.Owner, msg.Whom) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "%s is not following %s", msg.Owner, msg.Whom)
	}

	s.keeper.Unfollow(ctx, msg.Owner, msg.Whom)

	return &types.MsgUnfollowResponse{}, nil
}
