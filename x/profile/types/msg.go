package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetPrivate defines a SetPrivate message
type MsgSetPrivate struct {
	Owner   sdk.AccAddress `json:"owner"`
	Private []byte         `json:"private"`
}

// NewMsgSetSettings is a constructor function for MsgSetPrivate
func NewMsgSetPrivate(data []byte, owner sdk.AccAddress) MsgSetPrivate {
	return MsgSetPrivate{
		Private: data,
		Owner:   owner,
	}
}

// Route should return the name of the module
func (msg MsgSetPrivate) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetPrivate) Type() string { return "set_private" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPrivate) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Private) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Private cannot be empty")
	}
	if len(msg.Private) > 4*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Private cannot be greater than 4Kb")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPrivate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetPrivate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgSetPrivate defines a SetPrivate message
type MsgSetPublic struct {
	Owner  sdk.AccAddress `json:"owner"`
	Public Public         `json:"public"`
}

// NewMsgSetSettings is a constructor function for MsgSetPublic
func NewMsgSetPublic(public Public, owner sdk.AccAddress) MsgSetPublic {
	return MsgSetPublic{
		Public: public,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgSetPublic) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetPublic) Type() string { return "set_public" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPublic) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if msg.Public.Gender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Gender cannot be empty")
	}

	if !IsValidGender(string(msg.Public.Gender)) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Gender is not valid")
	}

	if !IsValidDate(msg.Public.Birthday) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Birthday is not valid")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPublic) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetPublic) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
