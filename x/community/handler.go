package community

import (
	"fmt"
	"time"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gofrs/uuid"
)

// NewHandler creates an sdk.Handler for all the community type messages
func NewHandler(cerberus cerberusapi.Cerberus, keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePost:
			return handleMsgCreatePost(ctx, keeper, msg)
		case MsgDeletePost:
			return handleMsgDeletePost(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreatePost(ctx sdk.Context, keeper Keeper, msg MsgCreatePost) (*sdk.Result, error) {
	id, _ := uuid.FromString(msg.UUID)
	keeper.CreatePost(ctx, msg.Owner, Post{
		UUID:          id,
		Owner:         msg.Owner,
		Title:         msg.Title,
		Category:      msg.Category,
		PreviewImage:  msg.PreviewImage,
		Text:          msg.Text,
		LikesCount:    0,
		DislikesCount: 0,
		CreatedAt:     time.Now().Unix(),
	})

	return &sdk.Result{}, nil
}

func handleMsgDeletePost(ctx sdk.Context, keeper Keeper, msg MsgDeletePost) (*sdk.Result, error) {
	id, _ := uuid.FromString(msg.UUID)

	keeper.DeletePost(ctx, msg.Owner, id)

	return &sdk.Result{}, nil
}
