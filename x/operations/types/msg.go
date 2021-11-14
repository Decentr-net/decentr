package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgDistributeRewards is a constructor function for MsgDistributeRewards
func NewMsgDistributeRewards(owner sdk.AccAddress, rewards []Reward) MsgDistributeRewards {
	return MsgDistributeRewards{
		Owner:   owner,
		Rewards: rewards,
	}
}

func NewReward(address sdk.AccAddress, reward sdk.Dec) Reward {
	return Reward{
		Receiver: address,
		Reward:   sdk.DecProto{Dec: reward},
	}
}

// Route should return the name of the module
func (m MsgDistributeRewards) Route() string { return RouterKey }

// Type should return the action
func (m MsgDistributeRewards) Type() string { return "distribute_rewards" }

// ValidateBasic runs stateless checks on the message
func (m MsgDistributeRewards) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Owner); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner")
	}

	if len(m.Rewards) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty rewards list")
	}

	if len(m.Rewards) > 1000 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "more than 1000 rewards")
	}

	for i, v := range m.Rewards {
		if err := sdk.VerifyAddressFormat(v.Receiver); err != nil {
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
	return []sdk.AccAddress{m.Owner}
}

func NewMsgResetAccount(owner, address sdk.AccAddress) MsgResetAccount {
	return MsgResetAccount{
		Owner:   owner,
		Address: address,
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
	return []sdk.AccAddress{m.Owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgResetAccount) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Owner); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner")
	}

	if err := sdk.VerifyAddressFormat(m.Address); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	return nil
}

func NewMsgBanAccount(owner, address sdk.AccAddress, ban bool) MsgBanAccount {
	return MsgBanAccount{
		Owner:   owner,
		Address: address,
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
	return []sdk.AccAddress{m.Owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgBanAccount) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Owner); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner")
	}

	if err := sdk.VerifyAddressFormat(m.Address); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	return nil
}

func NewMsgMint(owner sdk.AccAddress, coin sdk.Coin) MsgMint {
	return MsgMint{
		Owner: owner,
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
	return []sdk.AccAddress{m.Owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgMint) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Owner); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner")
	}

	if m.Coin.IsNil() || !m.Coin.IsValid() || m.Coin.IsZero() || m.Coin.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid coin")
	}
	return nil
}

func NewMsgBurn(owner sdk.AccAddress, coin sdk.Coin) MsgBurn {
	return MsgBurn{
		Owner: owner,
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
	return []sdk.AccAddress{m.Owner}
}

// ValidateBasic runs stateless checks on the message
func (m MsgBurn) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Owner); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner")
	}

	if m.Coin.IsNil() || !m.Coin.IsValid() || m.Coin.IsZero() || m.Coin.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid coin")
	}
	return nil
}
