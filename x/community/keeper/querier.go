package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gofrs/uuid"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/community/types"
)

const (
	QueryPopular = "popular"
	QueryPosts   = "posts"
	QueryUser    = "user"
)

type Post struct {
	UUID          string         `json:"uuid"`
	Owner         sdk.AccAddress `json:"owner"`
	Title         string         `json:"title"`
	PreviewImage  string         `json:"previewImage"`
	Category      types.Category `json:"category"`
	Text          string         `json:"text"`
	LikesCount    uint32         `json:"likesCount"`
	DislikesCount uint32         `json:"dislikesCount"`
	CreatedAt     int64          `json:"createdAt"`
}

// NewQuerier creates a new querier for community clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryPopular:
			return getIndexPostQuerier(popularityIndexBucket, getPopularityIndexKey)(ctx, path[1:], req, keeper)
		case QueryPosts:
			return getIndexPostQuerier(createdAtIndexBucket, getCreateAtIndexKey)(ctx, path[1:], req, keeper)
		case QueryUser:
			return queryUserPosts(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown community query endpoint")
		}
	}
}

// nolint: unparam
// queryPopular returns posts.
func queryUserPosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var (
		from  uuid.UUID
		limit = uint32(20)
	)

	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if path[1] != "" {
		from, err = uuid.FromString(path[1])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid from")
		}
	}

	if path[2] != "" {
		v, err := strconv.Atoi(path[2])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid limit")
		}
		limit = uint32(v)
	}

	p := keeper.ListUserPosts(ctx, owner, from, limit)

	res, err := codec.MarshalJSONIndent(keeper.cdc, postsToQuerierPosts(p))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func getIndexPostQuerier(index string, keyResolver func(p types.Post) []byte) func(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
		var (
			from     []byte
			category types.Category
			limit    = uint32(20)
		)

		if path[0] != "" || path[1] != "" {
			owner, err := sdk.AccAddressFromBech32(path[0])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}

			id, err := uuid.FromString(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}

			post := keeper.GetPostByKey(ctx, getPostKeeperKey(owner, id))
			from = keyResolver(post)
		}

		if path[2] != "" {
			v, err := strconv.Atoi(path[2])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid page")
			}
			category = types.Category(v)
		}

		if path[3] != "" {
			v, err := strconv.Atoi(path[3])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid page")
			}
			limit = uint32(v)
		}

		p, err := keeper.index.GetPosts(index, keeper.getPostResolver(ctx), category, from, limit)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, err.Error())
		}

		res, err := codec.MarshalJSONIndent(keeper.cdc, postsToQuerierPosts(p))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}

		return res, nil
	}
}

func postsToQuerierPosts(pp []types.Post) []Post {
	out := make([]Post, len(pp))

	for i, v := range pp {
		out[i] = Post{
			UUID:          v.UUID.String(),
			Owner:         v.Owner,
			Title:         v.Title,
			PreviewImage:  v.PreviewImage,
			Category:      v.Category,
			Text:          v.Text,
			LikesCount:    v.LikesCount,
			DislikesCount: v.DislikesCount,
			CreatedAt:     v.CreatedAt,
		}
	}

	return out
}
