package operations

import (
	"github.com/Decentr-net/decentr/x/operations/ante"
	"github.com/Decentr-net/decentr/x/operations/keeper"
	"github.com/Decentr-net/decentr/x/operations/types"
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

var (
	NewMinGasPriceDecorator = ante.NewMinGasPriceDecorator
)

type (
	Keeper               = keeper.Keeper
	MsgDistributeRewards = types.MsgDistributeRewards
	MsgResetAccount      = types.MsgResetAccount
	FixedGasParams       = types.FixedGasParams
)
