package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetLike{}, "community/SetLike", nil)
	cdc.RegisterConcrete(MsgCreatePost{}, "community/CreatePost", nil)
	cdc.RegisterConcrete(MsgDeletePost{}, "community/DeletePost", nil)
	cdc.RegisterConcrete(MsgFollow{}, "community/MsgFollow", nil)
	cdc.RegisterConcrete(MsgUnfollow{}, "community/MsgUnfollow", nil)
}

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
