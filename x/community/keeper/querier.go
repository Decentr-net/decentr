package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/gofrs/uuid"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/community/types"
)

const (
	QueryPost       = "post"
	QueryUser       = "user"
	QueryModerators = "moderators"
	QueryFollowees  = "followees"
)

const defaultLimit = 20

// NewQuerier creates a new querier for community clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryPost:
			return getPost(ctx, path[1:], req, keeper)
		case QueryUser:
			return queryUserPosts(ctx, path[1:], req, keeper)
		case QueryModerators:
			return queryModerators(ctx, keeper)
		case QueryFollowees:
			return queryFollowees(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown community query endpoint")
		}
	}
}

// nolint: unparam
// queryUserPosts returns posts for user.
func queryUserPosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	limit := defaultLimit
	from := uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	var err error

	if path[0] != "" {
		from, err = uuid.FromString(path[0])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid uuid")
		}
	}

	if path[1] != "" {
		var v int
		v, err = strconv.Atoi(path[1])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid limit")
		}
		limit = v
	}

	if from == uuid.Nil {
		// use max range

	}

	owner, err := sdk.AccAddressFromBech32(path[2])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address")
	}

	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), types.PostPrefix)

	it := store.ReverseIterator(getPostKeeperKey(owner, uuid.Nil), getPostKeeperKey(owner, from))
	defer it.Close()

	out := make([]types.Post, 0)
	for i := 0; i < limit && it.Valid(); it.Next() {
		var post types.Post
		keeper.cdc.MustUnmarshalBinaryBare(it.Value(), &post)
		out = append(out, post)

		i++
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, out)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func getPost(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address")
	}

	id, err := uuid.FromString(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid uuid")
	}

	p := keeper.GetPostByKey(ctx, getPostKeeperKey(owner, id))

	res, err := codec.MarshalJSONIndent(keeper.cdc, p)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryModerators(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	moderators := keeper.GetModerators(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, moderators)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryFollowees(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address")
	}

	followees := keeper.GetFollowees(ctx, owner)
	out := make([]string, len(followees))
	for idx, followee := range followees {
		out[idx] = followee.String()
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, out)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
