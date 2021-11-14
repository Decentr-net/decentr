package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(_ *codec.LegacyAmino) {
}

func RegisterInterfaces(_ cdctypes.InterfaceRegistry) {
}

var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
