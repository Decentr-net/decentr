package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Reward struct {
	Receiver sdk.AccAddress `json:"receiver"`
	ID       uint64         `json:"id"`
	Reward   uint64         `json:"reward"`
}

// MsgDistributeRewards defines a CreatePDV message
type MsgDistributeRewards struct {
	Owner   sdk.AccAddress `json:"owner"`
	Rewards []Reward       `json:"rewards"`
}

// NewMsgDistributeRewards is a constructor function for MsgDistributeRewards
func NewMsgDistributeRewards(owner sdk.AccAddress, rewards []Reward) MsgDistributeRewards {
	return MsgDistributeRewards{
		Owner:   owner,
		Rewards: rewards,
	}
}

// Route should return the name of the module
func (msg MsgDistributeRewards) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDistributeRewards) Type() string { return "distribute_rewards" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDistributeRewards) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Owner is empty")
	}

	if len(msg.Rewards) > 1000 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Number of rewards can't be greater than 1000")
	}

	for _, reward := range msg.Rewards {
		if reward.Receiver.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Receiver is empty")
		}

		if reward.Reward == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Reward can't be zero")
		}

		if reward.ID == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID can't be zero")
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDistributeRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDistributeRewards) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgResetAccount struct {
	Owner        sdk.AccAddress `json:"owner"`
	AccountOwner sdk.AccAddress `json:"accountOwner"`
}

func NewMsgResetAccount(owner, accountOwner sdk.AccAddress) MsgResetAccount {
	return MsgResetAccount{
		Owner:        owner,
		AccountOwner: accountOwner,
	}
}

// Route should return the name of the module
func (msg MsgResetAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgResetAccount) Type() string { return "reset_account" }

// GetSignBytes encodes the message for signing
func (msg MsgResetAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgResetAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic runs stateless checks on the message
func (msg MsgResetAccount) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Owner is empty")
	}
	if msg.AccountOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "AccountOwner is empty")
	}
	return nil
}
