package types

import (
	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreatePDV defines a CreatePDV message
type MsgCreatePDV struct {
	Timestamp uint64         `json:"timestamp"`
	Address   string         `json:"address"`
	Owner     sdk.AccAddress `json:"owner"`
	DataType  PDVType        `json:"type"`
}

// NewMsgCreatePDV is a constructor function for MsgCreatePDV
func NewMsgCreatePDV(timestamp uint64, value string, dataType PDVType, owner sdk.AccAddress) MsgCreatePDV {
	return MsgCreatePDV{
		Timestamp: timestamp,
		Address:   value,
		Owner:     owner,
		DataType:  dataType,
	}
}

// Route should return the name of the module
func (msg MsgCreatePDV) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreatePDV) Type() string { return "create_pdv" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreatePDV) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if msg.DataType < PDVTypeCookie || msg.DataType > PDVTypeLoginCookie {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid type")
	}
	if !cerberusapi.IsAddressValid(msg.Address) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address")
	}
	if msg.Timestamp == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Timestamp can not be 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreatePDV) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreatePDV) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
