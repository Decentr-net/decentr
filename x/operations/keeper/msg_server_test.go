package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/operations/types"
)

func TestMsgServer_DistributeRewards(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	owner, addr1, addr2 := NewAccAddress(), NewAccAddress(), NewAccAddress()
	set.keeper.SetParams(ctx, types.Params{
		Supervisors: []string{owner.String()},
		FixedGas:    types.DefaultFixedGasParams(),
		MinGasPrice: types.DefaultMinGasPrice,
	})

	_, err := s.DistributeRewards(sdk.WrapSDKContext(ctx), &types.MsgDistributeRewards{
		Owner: owner.String(),
		Rewards: []types.Reward{
			{Receiver: addr1.String(), Reward: 1},
			{Receiver: addr2.String(), Reward: 2},
		},
	})
	require.NoError(t, err)

	// check token
}

func TestMsgServer_DistributeRewards_Unauthorized(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	owner, addr1 := NewAccAddress(), NewAccAddress()

	_, err := s.DistributeRewards(sdk.WrapSDKContext(ctx), &types.MsgDistributeRewards{
		Owner: owner.String(),
		Rewards: []types.Reward{
			{Receiver: addr1.String(), Reward: 1},
		},
	})
	require.Error(t, err)
	require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrUnauthorized))
}

func TestMsgServer_ResetAccount(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	owner, address := NewAccAddress(), NewAccAddress()
	set.keeper.SetParams(ctx, types.Params{
		Supervisors: []string{owner.String()},
		FixedGas:    types.DefaultFixedGasParams(),
		MinGasPrice: types.DefaultMinGasPrice,
	})

	_, err := s.ResetAccount(sdk.WrapSDKContext(ctx), &types.MsgResetAccount{
		Owner:   owner.String(),
		Address: address.String(),
	})
	require.NoError(t, err)
	// check token
	// check blog
	// check likes
}

func TestMsgServer_ResetAccount_SelfReset(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	address := NewAccAddress()
	_, err := s.ResetAccount(sdk.WrapSDKContext(ctx), &types.MsgResetAccount{
		Owner:   address.String(),
		Address: address.String(),
	})
	require.NoError(t, err)
	// check token
	// check blog
	// check likes
}

func TestMsgServer_ResetAccount_Unauthorized(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	owner, address := NewAccAddress(), NewAccAddress()

	_, err := s.ResetAccount(sdk.WrapSDKContext(ctx), &types.MsgResetAccount{
		Owner:   owner.String(),
		Address: address.String(),
	})
	require.Error(t, err)
	require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrUnauthorized))
}

func TestMsgServer_BanAccount(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	owner, address := NewAccAddress(), NewAccAddress()
	set.keeper.SetParams(ctx, types.Params{
		Supervisors: []string{owner.String()},
		FixedGas:    types.DefaultFixedGasParams(),
		MinGasPrice: types.DefaultMinGasPrice,
	})

	_, err := s.BanAccount(sdk.WrapSDKContext(ctx), &types.MsgBanAccount{
		Owner:   owner.String(),
		Address: address.String(),
		Ban:     true,
	})
	require.NoError(t, err)
	require.True(t, set.tokenKeeper.IsBanned(ctx, address))

	_, err = s.BanAccount(sdk.WrapSDKContext(ctx), &types.MsgBanAccount{
		Owner:   owner.String(),
		Address: address.String(),
		Ban:     false,
	})
	require.NoError(t, err)
	require.False(t, set.tokenKeeper.IsBanned(ctx, address))
}

func TestMsgServer_BanAccount_Unauthorized(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	owner, address := NewAccAddress(), NewAccAddress()

	_, err := s.BanAccount(sdk.WrapSDKContext(ctx), &types.MsgBanAccount{
		Owner:   owner.String(),
		Address: address.String(),
		Ban:     true,
	})
	require.Error(t, err)
	require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrUnauthorized))
}

func TestMsgServer_Mint(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)
	bk := set.bankKeeper.(bankkeeper.Keeper)

	owner := NewAccAddress()
	set.keeper.SetParams(ctx, types.Params{
		Supervisors: []string{owner.String()},
		FixedGas:    types.DefaultFixedGasParams(),
		MinGasPrice: types.DefaultMinGasPrice,
	})

	coin := sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000))
	_, err := s.Mint(sdk.WrapSDKContext(ctx), &types.MsgMint{
		Owner: owner.String(),
		Coin:  coin,
	})
	require.NoError(t, err)

	require.Equal(t, coin, bk.GetSupply(ctx, config.DefaultBondDenom))
	require.Equal(t, coin, bk.GetBalance(ctx, owner, config.DefaultBondDenom))
}

func TestMsgServer_Mint_Unauthorized(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	_, err := s.Mint(sdk.WrapSDKContext(ctx), &types.MsgMint{
		Owner: NewAccAddress().String(),
		Coin:  sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000)),
	})
	require.Error(t, err)
	require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrUnauthorized))
}

func TestMsgServer_Burn(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)
	bk := set.bankKeeper.(bankkeeper.Keeper)

	owner := NewAccAddress()
	coin := sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000))

	require.NoError(t, bk.MintCoins(ctx, types.ModuleName, sdk.Coins{coin}))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.Coins{coin}))

	set.keeper.SetParams(ctx, types.Params{
		Supervisors: []string{owner.String()},
		FixedGas:    types.DefaultFixedGasParams(),
		MinGasPrice: types.DefaultMinGasPrice,
	})

	_, err := s.Burn(sdk.WrapSDKContext(ctx), &types.MsgBurn{
		Owner: owner.String(),
		Coin:  coin,
	})
	require.NoError(t, err)

	require.Equal(t, sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(0)), bk.GetSupply(ctx, config.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(0)), bk.GetBalance(ctx, owner, config.DefaultBondDenom))
}

func TestMsgServer_Burn_Unauthorized(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewMsgServer(set.keeper, set.bankKeeper, set.tokenKeeper, set.communityKeeper)

	_, err := s.Burn(sdk.WrapSDKContext(ctx), &types.MsgBurn{
		Owner: NewAccAddress().String(),
		Coin:  sdk.NewCoin(config.DefaultBondDenom, sdk.NewInt(1000)),
	})
	require.Error(t, err)
	require.True(t, sdkerrors.IsOf(err, sdkerrors.ErrUnauthorized))
}
