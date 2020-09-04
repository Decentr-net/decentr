package pdv

import (
	"github.com/Decentr-net/decentr/x/pdv/keeper"
	"github.com/Decentr-net/decentr/x/pdv/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	FlagCerberusAddr = types.FlagCerberusAddr
)

var (
	NewKeeper     = keeper.NewKeeper
	NewStats      = keeper.NewStats
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper       = keeper.Keeper
	Stats        = keeper.Stats
	PDV          = types.PDV
	MsgCreatePDV = types.MsgCreatePDV
)
