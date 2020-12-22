package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gofrs/uuid"
)

// NewHandler creates an sdk.Handler for all the community type messages
func NewHandler(keeper Keeper, moderatorAddress sdk.AccAddress) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePost:
			return handleMsgCreatePost(ctx, keeper, msg)
		case MsgDeletePost:
			return handleMsgDeletePost(ctx, keeper, msg, moderatorAddress)
		case MsgSetLike:
			return handleMsgSetLike(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreatePost(ctx sdk.Context, keeper Keeper, msg MsgCreatePost) (*sdk.Result, error) {
	id, _ := uuid.FromString(msg.UUID)
	keeper.CreatePost(ctx, Post{
		UUID:          id,
		Owner:         msg.Owner,
		Title:         msg.Title,
		Category:      msg.Category,
		PreviewImage:  msg.PreviewImage,
		Text:          msg.Text,
		LikesCount:    0,
		DislikesCount: 0,
		PDV:           sdk.NewInt(0),
		CreatedAt:     uint64(ctx.BlockTime().Unix()),
	})

	return &sdk.Result{}, nil
}

func handleMsgDeletePost(ctx sdk.Context, keeper Keeper, msg MsgDeletePost, moderatorAddress sdk.AccAddress) (*sdk.Result, error) {
	if !msg.Owner.Equals(moderatorAddress) && !msg.Owner.Equals(msg.PostOwner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	postUUID, _ := uuid.FromString(msg.PostUUID)
	keeper.DeletePost(ctx, msg.PostOwner, postUUID)

	if msg.Owner.Equals(moderatorAddress) {
		ctx.Logger().Info("moderator deleted post %s %s", msg.PostOwner, msg.PostUUID)
	}

	return &sdk.Result{}, nil
}

func handleMsgSetLike(ctx sdk.Context, keeper Keeper, msg MsgSetLike) (*sdk.Result, error) {
	postUUID, _ := uuid.FromString(msg.PostUUID)
	keeper.SetLike(ctx, Like{
		PostOwner: msg.PostOwner,
		PostUUID:  postUUID,
		Owner:     msg.Owner,
		Weight:    msg.Weight,
	})
	return &sdk.Result{}, nil
}
