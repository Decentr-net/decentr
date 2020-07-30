package profile

import (
	"github.com/Decentr-net/decentr/x/profile/keeper"
	"github.com/Decentr-net/decentr/x/profile/types"
)

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper        = keeper.Keeper
	Profile       = types.Profile
	MsgSetPrivate = types.MsgSetPrivate
	MsgSetPublic  = types.MsgSetPublic
)
