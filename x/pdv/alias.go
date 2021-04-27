package pdv

import (
	"github.com/Decentr-net/decentr/x/pdv/ante"
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

var (
	NewCreateAccountDecorator = ante.NewCreateAccountDecorator
)

type (
	Keeper                 = keeper.Keeper
	MsgDistributeRewards   = types.MsgDistributeRewards
	MsgResetAccount        = types.MsgResetAccount
	CreateAccountDecorator = ante.CreateAccountDecorator
	FixedGasParams         = types.FixedGasParams
)
