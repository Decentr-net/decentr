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

func NewReward(address sdk.AccAddress, reward uint64) Reward {
	return Reward{
		Receiver: address.String(),
		Reward:   reward,
	}
}

// Route should return the name of the module
func (m MsgDistributeRewards) Route() string { return RouterKey }

// Type should return the action
func (m MsgDistributeRewards) Type() string { return "distribute_rewards" }

// ValidateBasic runs stateless checks on the message
func (m MsgDistributeRewards) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("invalid owner address")
	}

	if len(m.Rewards) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty rewards list")
	}

	if len(m.Rewards) > 1000 {
		return sdkerrors.ErrInvalidRequest.Wrap("more than 1000 rewards")
	}

	for _, reward := range m.Rewards {
		if _, err := sdk.AccAddressFromBech32(reward.Receiver); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
		}

		if reward.Reward == 0 {
			return sdkerrors.ErrInvalidRequest.Wrapf("zero reward for %s", reward.Receiver)
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
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
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
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgResetAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", err)
	}

	return nil
}

func NewMsgBanAccount(owner, address sdk.AccAddress, ban bool) MsgBanAccount {
	return MsgBanAccount{
		Owner:   owner.String(),
		Address: address.String(),
		Ban:     ban,
	}
}

// Route should return the name of the module
func (m MsgBanAccount) Route() string { return RouterKey }

// Type should return the action
func (m MsgBanAccount) Type() string { return "ban_account" }

// GetSignBytes encodes the message for signing
func (m MsgBanAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgBanAccount) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgBanAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", err)
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
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
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
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgBurn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}
	if m.Coin.IsNil() || !m.Coin.IsValid() || m.Coin.IsZero() || m.Coin.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid coin")
	}
	return nil
}
