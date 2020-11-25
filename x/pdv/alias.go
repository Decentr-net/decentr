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
	NewIndex      = keeper.NewIndex
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper       = keeper.Keeper
	Index        = keeper.Index
	PDV          = types.PDV
	MsgCreatePDV = types.MsgCreatePDV
)
