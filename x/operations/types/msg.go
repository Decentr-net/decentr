package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgDistributeRewards is a constructor function for MsgDistributeRewards
func NewMsgDistributeRewards(owner sdk.AccAddress, rewards []Reward) MsgDistributeRewards {
	return MsgDistributeRewards{
		Owner:   owner.String(),
		Rewards: rewards,
	}
}

func NewReward(address sdk.AccAddress, reward sdk.Dec) Reward {
	return Reward{
		Receiver: address.String(),
		Reward:   sdk.DecProto{Dec: reward},
	}
}

// Route should return the name of the module
func (m MsgDistributeRewards) Route() string { return RouterKey }

// Type should return the action
func (m MsgDistributeRewards) Type() string { return "distribute_rewards" }

// ValidateBasic runs stateless checks on the message
func (m MsgDistributeRewards) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner")
	}

	if len(m.Rewards) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty rewards list")
	}

	if len(m.Rewards) > 1000 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "more than 1000 rewards")
	}

	for i, v := range m.Rewards {
		receiver, err := sdk.AccAddressFromBech32(v.Receiver)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid reward %d: invalid receiver", i+1)
		}

		if err := sdk.VerifyAddressFormat(receiver); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid reward %d: invalid receiver", i+1)
		}

		reward := v.Reward.Dec
		if reward.IsNil() || reward.IsZero() || reward.IsNegative() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid reward %d: invalid reward", i+1)
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgDistributeRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgDistributeRewards) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

func NewMsgResetAccount(owner, address sdk.AccAddress) MsgResetAccount {
	return MsgResetAccount{
		Owner:   owner.String(),
		Address: address.String(),
	}
}

// Route should return the name of the module
func (m MsgResetAccount) Route() string { return RouterKey }

// Type should return the action
func (m MsgResetAccount) Type() string { return "reset_account" }

// GetSignBytes encodes the message for signing
func (m MsgResetAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgResetAccount) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message
func (m MsgResetAccount) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}
	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner: %s", err)
	}

	address, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address address: %s", err)
	}
	if err := sdk.VerifyAddressFormat(address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", err)
	}

	return nil
}

func NewMsgMint(owner sdk.AccAddress, coin sdk.Coin) MsgMint {
	return MsgMint{
		Owner: owner.String(),
		Coin:  coin,
	}
}

// Route should return the name of the module
func (m MsgMint) Route() string { return RouterKey }

// Type should return the action
func (m MsgMint) Type() string { return "mint" }

// GetSignBytes encodes the message for signing
func (m MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgMint) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message
func (m MsgMint) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner: %s", err)
	}

	if m.Coin.IsNil() || !m.Coin.IsValid() || m.Coin.IsZero() || m.Coin.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid coin")
	}
	return nil
}

func NewMsgBurn(owner sdk.AccAddress, coin sdk.Coin) MsgBurn {
	return MsgBurn{
		Owner: owner.String(),
		Coin:  coin,
	}
}

// Route should return the name of the module
func (m MsgBurn) Route() string { return RouterKey }

// Type should return the action
func (m MsgBurn) Type() string { return "burn" }

// GetSignBytes encodes the message for signing
func (m MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgBurn) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message
func (m MsgBurn) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}
	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner: %s", err)
	}

	if m.Coin.IsNil() || !m.Coin.IsValid() || m.Coin.IsZero() || m.Coin.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid coin")
	}
	return nil
}
