package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreatePDV defines a CreatePDV message
type MsgCreatePDV struct {
	Owner sdk.AccAddress `json:"owner"`
	ID    uint64         `json:"id"`
}

// NewMsgCreatePDV is a constructor function for MsgCreatePDV
func NewMsgCreatePDV(owner sdk.AccAddress, id uint64) MsgCreatePDV {
	return MsgCreatePDV{
		Owner: owner,
		ID:    id,
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

	if msg.ID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id can't be zero")
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
