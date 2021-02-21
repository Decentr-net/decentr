package pdv

import (
	"github.com/Decentr-net/decentr/x/pdv/keeper"
	"github.com/Decentr-net/decentr/x/pdv/types"
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
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper               = keeper.Keeper
	MsgDistributeRewards = types.MsgDistributeRewards
)
