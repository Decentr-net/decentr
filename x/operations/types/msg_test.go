package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
)

func TestMsgDistributeRewards_ValidateBasic(t *testing.T) {
	valid := NewMsgDistributeRewards(NewAccAddress(), []Reward{
		NewReward(NewAccAddress(), sdk.NewDec(10)),
	})

	alter := func(f func(m *MsgDistributeRewards)) MsgDistributeRewards {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		for i := 0; i < 12; i++ {
			m.Rewards = append(m.Rewards, m.Rewards...)
		}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{
			{Receiver: nil, Reward: sdk.DecProto{sdk.OneDec()}},
		}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{
			NewReward(NewAccAddress(), sdk.ZeroDec()),
		}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{
			NewReward(NewAccAddress(), sdk.OneDec()),
			NewReward(NewAccAddress(), sdk.ZeroDec()),
		}
	}).ValidateBasic())
}

func TestMsgResetAccount_ValidateBasic(t *testing.T) {
	valid := NewMsgResetAccount(NewAccAddress(), NewAccAddress())

	alter := func(f func(m *MsgResetAccount)) MsgResetAccount {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgResetAccount) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgResetAccount) {
		m.Address = nil
	}).ValidateBasic())
}

func TestMsgMint_ValidateBasic(t *testing.T) {
	valid := NewMsgMint(NewAccAddress(), sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000)))

	alter := func(f func(m *MsgMint)) MsgMint {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgMint) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgMint) {
		m.Coin.Denom = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgMint) {
		m.Coin.Amount = sdk.Int{}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgMint) {
		m.Coin.Amount = sdk.NewInt(-1)
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgMint) {
		m.Coin.Amount = sdk.ZeroInt()
	}).ValidateBasic())
}

func TestMsgBurn_ValidateBasic(t *testing.T) {
	valid := NewMsgBurn(NewAccAddress(), sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000)))

	alter := func(f func(m *MsgBurn)) MsgBurn {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgBurn) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBurn) {
		m.Coin.Denom = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBurn) {
		m.Coin.Amount = sdk.Int{}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBurn) {
		m.Coin.Amount = sdk.NewInt(-1)
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBurn) {
		m.Coin.Amount = sdk.ZeroInt()
	}).ValidateBasic())
}
