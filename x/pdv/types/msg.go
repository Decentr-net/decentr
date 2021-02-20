package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreatePDV defines a CreatePDV message
type MsgCreatePDV struct {
	Owner    sdk.AccAddress `json:"owner"`
	Receiver sdk.AccAddress `json:"receiver"`
	ID       uint64         `json:"id"`
	Reward   uint64         `json:"reward"`
}

// NewMsgCreatePDV is a constructor function for MsgCreatePDV
func NewMsgCreatePDV(owner sdk.AccAddress, id uint64, receiver sdk.AccAddress, reward uint64) MsgCreatePDV {
	return MsgCreatePDV{
		Owner:    owner,
		Receiver: receiver,
		ID:       id,
		Reward:   reward,
	}
}

// Route should return the name of the module
func (msg MsgCreatePDV) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreatePDV) Type() string { return "create_pdv" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreatePDV) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Owner is empty")
	}

	if msg.Receiver.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Receiver is empty")
	}

	if msg.Reward == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Reward can't be zero")
	}

	if msg.ID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID can't be zero")
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
