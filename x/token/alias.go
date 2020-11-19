package token

import (
	"github.com/Decentr-net/decentr/x/token/keeper"
	"github.com/Decentr-net/decentr/x/token/types"
)

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewStats      = keeper.NewStats
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper
	Stats  = keeper.Stats
)
