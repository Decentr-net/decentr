package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDistributeRewards{}, "operations/DistributeRewards", nil)
	cdc.RegisterConcrete(MsgResetAccount{}, "operations/MsgResetAccount", nil)
}

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
