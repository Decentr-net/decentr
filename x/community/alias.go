package community

import (
	"github.com/Decentr-net/decentr/x/community/keeper"
	"github.com/Decentr-net/decentr/x/community/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	FlagModeratorAddr = types.FlagModeratorAddr
)

var (
	NewKeeper     = keeper.NewKeeper
	NewIndex      = keeper.NewIndex
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper        = keeper.Keeper
	Index         = keeper.Index
	Post          = types.Post
	Like          = types.Like
	MsgCreatePost = types.MsgCreatePost
	MsgDeletePost = types.MsgDeletePost
	MsgSetLike    = types.MsgSetLike
)
