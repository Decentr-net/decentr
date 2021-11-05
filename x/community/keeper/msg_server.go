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

	owner, _ := sdk.AccAddressFromBech32(msg.Post.Owner)
	id, _ := uuid.FromString(msg.Post.Uuid)

	if p := s.keeper.GetPost(ctx, postKey(owner, id)); p != (types.Post{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrConflict, "post %s already exists", msg.Post.Address())
	}

	s.keeper.SetPost(ctx, msg.Post)

	return &types.MsgCreatePostResponse{}, nil
}

func (s msgServer) DeletePost(goCtx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Owner != msg.PostOwner {
		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		if !s.keeper.IsModerator(ctx, owner) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a moderator or post owner", msg.Owner)
		}
	}

	owner, _ := sdk.AccAddressFromBech32(msg.PostOwner)
	id, _ := uuid.FromString(msg.PostUuid)
	key := postKey(owner, id)

	if p := s.keeper.GetPost(ctx, key); p == (types.Post{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "post %s/%s not found", msg.PostOwner, msg.PostUuid)
	}

	s.keeper.DeletePost(ctx, key)

	return &types.MsgDeletePostResponse{}, nil
}

func (s msgServer) SetLike(goCtx context.Context, msg *types.MsgSetLike) (*types.MsgSetLikeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	postOwner, _ := sdk.AccAddressFromBech32(msg.Like.PostOwner)
	postUUID, _ := uuid.FromString(msg.Like.PostUuid)
	owner, _ := sdk.AccAddressFromBech32(msg.Like.Owner)

	if owner.Equals(postOwner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "self-like is disabled")
	}

	if p := s.keeper.GetPost(ctx, postKey(postOwner, postUUID)); p == (types.Post{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "post %s/%s not found", msg.Like.PostOwner, msg.Like.PostUuid)
	}

	diff := msg.Like.Weight - s.keeper.GetLike(ctx, likeKey(postKey(postOwner, postUUID), owner)).Weight
	s.tokenKeeper.IncTokens(ctx, postOwner, sdk.NewDec(int64(diff)))

	s.keeper.SetLike(ctx, msg.Like)

	return &types.MsgSetLikeResponse{}, nil
}

func (s msgServer) Follow(goCtx context.Context, msg *types.MsgFollow) (*types.MsgFollowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	whom, _ := sdk.AccAddressFromBech32(msg.Whom)

	if owner.Equals(whom) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "self-follow is disabled")
	}

	if s.keeper.IsFollowed(ctx, owner, whom) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrConflict, "%s already follows %s", owner, whom)
	}

	s.keeper.Follow(ctx, owner, whom)

	return &types.MsgFollowResponse{}, nil
}

func (s msgServer) Unfollow(goCtx context.Context, msg *types.MsgUnfollow) (*types.MsgUnfollowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	whom, _ := sdk.AccAddressFromBech32(msg.Whom)

	if !s.keeper.IsFollowed(ctx, owner, whom) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "%s is not following %s", owner, whom)
	}

	s.keeper.Unfollow(ctx, owner, whom)

	return &types.MsgUnfollowResponse{}, nil
}
