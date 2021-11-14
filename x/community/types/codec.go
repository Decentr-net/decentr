package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreatePost{}, "community/MsgCreatePost", nil)
	cdc.RegisterConcrete(MsgDeletePost{}, "community/MsgDeletePost", nil)
	cdc.RegisterConcrete(MsgSetLike{}, "community/MsgSetLike", nil)
	cdc.RegisterConcrete(MsgFollow{}, "community/MsgFollow", nil)
	cdc.RegisterConcrete(MsgUnfollow{}, "community/MsgUnfollow", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePost{},
		&MsgDeletePost{},
		&MsgSetLike{},
		&MsgFollow{},
		&MsgUnfollow{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
