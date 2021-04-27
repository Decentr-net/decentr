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
	DefaultParamspace = types.DefaultParamspace
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper         = keeper.Keeper
	Post           = types.Post
	Like           = types.Like
	MsgCreatePost  = types.MsgCreatePost
	MsgDeletePost  = types.MsgDeletePost
	MsgFollow      = types.MsgFollow
	MsgUnfollow    = types.MsgUnfollow
	MsgSetLike     = types.MsgSetLike
	FixedGasParams = types.FixedGasParams
)
