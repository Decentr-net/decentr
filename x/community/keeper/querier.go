package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gofrs/uuid"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/community/types"
	"github.com/Decentr-net/decentr/x/utils"
)

const (
	QueryPopular       = "popular"
	QueryPosts         = "posts"
	QueryPost          = "post"
	QueryUser          = "user"
	QueryLikedPosts    = "liked-posts"
	QueryModeratorAddr = "moderator-addr"
	QueryModerators    = "moderators"
)

const defaultLimit = 20

type Post struct {
	UUID          string         `json:"uuid"`
	Owner         sdk.AccAddress `json:"owner"`
	Title         string         `json:"title"`
	PreviewImage  string         `json:"previewImage"`
	Category      types.Category `json:"category"`
	Text          string         `json:"text"`
	LikesCount    uint32         `json:"likesCount"`
	DislikesCount uint32         `json:"dislikesCount"`
	CreatedAt     uint64         `json:"createdAt"`
	PDV           float64        `json:"pdv" amino:"unsafe"`
}

// NewQuerier creates a new querier for community clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryPopular:
			return getPopularPosts(ctx, path[1:], req, keeper)
		case QueryPost:
			return getPost(ctx, path[1:], req, keeper)
		case QueryPosts:
			return getRecentPosts(ctx, path[1:], req, keeper)
		case QueryUser:
			return queryUserPosts(ctx, path[1:], req, keeper)
		case QueryLikedPosts:
			return queryUserLikedPosts(ctx, path[1:], req, keeper)
		case QueryModeratorAddr:
			return queryModeratorAddr(ctx, keeper)
		case QueryModerators:
			return queryModerators(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown community query endpoint")
		}
	}
}

// nolint: unparam
// queryPopular returns posts.
func queryUserPosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, from, limit, err := extractCommonGetParameters(path)
	if err != nil {
		return nil, err
	}

	p := keeper.ListUserPosts(ctx, owner, from, limit)

	res, err := codec.MarshalJSONIndent(keeper.cdc, postsToQuerierPosts(p))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
// queryPopular returns map with post and its like weight.
func queryUserLikedPosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetUserLikedPosts(ctx, owner))
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

	res, err := codec.MarshalJSONIndent(keeper.cdc, Post{
		UUID:          p.UUID.String(),
		Owner:         p.Owner,
		Title:         p.Title,
		PreviewImage:  p.PreviewImage,
		Category:      p.Category,
		Text:          p.Text,
		LikesCount:    p.LikesCount,
		DislikesCount: p.DislikesCount,
		CreatedAt:     p.CreatedAt,
		PDV:           utils.TokenToFloat64(p.PDV),
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func getRecentPosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, id, limit, err := extractCommonGetParameters(path)
	if err != nil {
		return nil, err
	}

	var from []byte
	if !owner.Empty() && id != uuid.Nil {
		from = getCreateAtIndexKey(keeper.GetPostByKey(ctx, getPostKeeperKey(owner, id)))
	}

	var category types.Category
	if path[3] != "" {
		v, err := strconv.Atoi(path[3])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid category")
		}
		category = types.Category(v)
		if category < types.UndefinedCategory || category > types.SportsCategory {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown category")
		}
	}

	p := keeper.GetRecentPosts(ctx, category, from, limit)
	res, err := codec.MarshalJSONIndent(keeper.cdc, postsToQuerierPosts(p))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func getPopularPosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, id, limit, err := extractCommonGetParameters(path)
	if err != nil {
		return nil, err
	}

	var from []byte
	if !owner.Empty() && id != uuid.Nil {
		from = getPopularityIndexKey(keeper.GetPostByKey(ctx, getPostKeeperKey(owner, id)))
	}

	var category types.Category
	if path[3] != "" {
		v, err := strconv.Atoi(path[3])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid category")
		}
		category = types.Category(v)
	}

	var interval Interval
	if path[4] != "" {
		v, err := strconv.Atoi(path[4])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid interval")
		}
		interval = Interval(v)
	}
	if interval == InvalidInterval || interval > MonthInterval {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid interval")
	}

	p := keeper.GetPopularPosts(ctx, interval, category, from, limit)
	res, err := codec.MarshalJSONIndent(keeper.cdc, postsToQuerierPosts(p))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func postsToQuerierPosts(pp []types.Post) []Post {
	out := make([]Post, 0, len(pp))

	for _, v := range pp {
		if v.UUID == uuid.Nil {
			continue
		}

		out = append(out, Post{
			UUID:          v.UUID.String(),
			Owner:         v.Owner,
			Title:         v.Title,
			PreviewImage:  v.PreviewImage,
			Category:      v.Category,
			Text:          v.Text,
			LikesCount:    v.LikesCount,
			DislikesCount: v.DislikesCount,
			CreatedAt:     v.CreatedAt,
			PDV:           utils.TokenToFloat64(v.PDV),
		})
	}

	return out
}

func extractCommonGetParameters(path []string) (owner sdk.AccAddress, id uuid.UUID, limit uint32, err error) {
	limit = defaultLimit

	if path[0] != "" {
		owner, err = sdk.AccAddressFromBech32(path[0])
		if err != nil {
			err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address")
			return
		}
	}

	if path[1] != "" {
		id, err = uuid.FromString(path[1])
		if err != nil {
			err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid uuid")
			return
		}
	}

	if path[2] != "" {
		var v int
		v, err = strconv.Atoi(path[2])
		if err != nil {
			err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid limit")
			return
		}
		limit = uint32(v)
	}

	return
}

func queryModeratorAddr(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	moderators := keeper.GetModerators(ctx)
	if len(moderators) != 0 {
		return []byte(moderators[0]), nil
	}

	return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, "moderator not defined")
}

func queryModerators(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	moderators := keeper.GetModerators(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, moderators)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
