package decentr

import (
	"github.com/Decentr-net/decentr/x/decentr/keeper"
	"github.com/Decentr-net/decentr/x/decentr/types"
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
	Keeper       = keeper.Keeper
	PDV          = types.PDV
	MsgCreatePDV = types.MsgCreatePDV
)
