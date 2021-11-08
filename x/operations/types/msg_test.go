package types

import (
	"testing"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDistributeRewards_ValidateBasic(t *testing.T) {
	valid := NewMsgDistributeRewards(NewAccAddress(), []Reward{
		NewReward(NewAccAddress(), 10),
	})

	alter := func(f func(m *MsgDistributeRewards)) MsgDistributeRewards {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Owner = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Owner = "1"
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
			{Receiver: "", Reward: 1},
		}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{
			{Receiver: "1", Reward: 1},
		}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{
			NewReward(NewAccAddress(), 0),
		}
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDistributeRewards) {
		m.Rewards = []Reward{
			NewReward(NewAccAddress(), 1),
			NewReward(NewAccAddress(), 0),
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
		m.Owner = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgResetAccount) {
		m.Address = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgResetAccount) {
		m.Owner = "1"
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgResetAccount) {
		m.Address = "1"
	}).ValidateBasic())
}

func TestMsgBanAccount_ValidateBasic(t *testing.T) {
	valid := NewMsgBanAccount(NewAccAddress(), NewAccAddress(), true)

	alter := func(f func(m *MsgBanAccount)) MsgBanAccount {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgBanAccount) {
		m.Owner = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBanAccount) {
		m.Address = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBanAccount) {
		m.Owner = "1"
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBanAccount) {
		m.Address = "1"
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
		m.Owner = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgMint) {
		m.Owner = "1"
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
		m.Owner = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgBurn) {
		m.Owner = "1"
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
