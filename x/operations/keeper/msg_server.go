package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/operations/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper          Keeper
	bankKeeper      types.BankKeeper
	tokenKeeper     types.TokenKeeper
	communityKeeper types.CommunityKeeper
}

// NewMsgServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServer(
	keeper Keeper,
	bankKeeper types.BankKeeper,
	tokenKeeper types.TokenKeeper,
	communityKeeper types.CommunityKeeper,
) types.MsgServer {
	return &msgServer{
		keeper:          keeper,
		bankKeeper:      bankKeeper,
		tokenKeeper:     tokenKeeper,
		communityKeeper: communityKeeper,
	}
}

func (s msgServer) DistributeRewards(goCtx context.Context, msg *types.MsgDistributeRewards) (*types.MsgDistributeRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)

	if !s.keeper.IsSupervisor(ctx, owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a supervisor", msg.Owner)
	}

	for _, v := range msg.Rewards {
		address, _ := sdk.AccAddressFromBech32(v.Receiver)
		s.tokenKeeper.IncTokens(ctx, address, v.Reward.Dec)
	}

	return &types.MsgDistributeRewardsResponse{}, nil
}

func (s msgServer) ResetAccount(goCtx context.Context, msg *types.MsgResetAccount) (*types.MsgResetAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	address, _ := sdk.AccAddressFromBech32(msg.Address)

	if !s.keeper.IsSupervisor(ctx, owner) && !owner.Equals(address) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not an owner or supervisor", msg.Owner)
	}

	// reset account in other modules
	s.tokenKeeper.ResetAccount(ctx, address)
	s.communityKeeper.ResetAccount(ctx, address)

	return &types.MsgResetAccountResponse{}, nil
}

func (s msgServer) BanAccount(goCtx context.Context, msg *types.MsgBanAccount) (*types.MsgBanAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	if !s.keeper.IsSupervisor(ctx, owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a supervisor", msg.Owner)
	}

	address, _ := sdk.AccAddressFromBech32(msg.Address)
	s.tokenKeeper.SetBan(ctx, address, msg.Ban)

	return &types.MsgBanAccountResponse{}, nil
}

func (s msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	if !s.keeper.IsSupervisor(ctx, owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a supervisor", msg.Owner)
	}

	// mint new tokens and send it to owner
	if err := s.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(msg.Coin)); err != nil {
		return nil, err
	}
	if err := s.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, owner, sdk.NewCoins(msg.Coin)); err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}

func (s msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	if !s.keeper.IsSupervisor(ctx, owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a supervisor", msg.Owner)
	}

	// send tokens to module and burn it
	if err := s.bankKeeper.SendCoinsFromAccountToModule(
		ctx, owner, types.ModuleName, sdk.NewCoins(msg.Coin)); err != nil {
		return nil, err
	}
	if err := s.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(msg.Coin)); err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}
