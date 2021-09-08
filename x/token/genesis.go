package token

import (
	"github.com/Decentr-net/decentr/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params   types.Params                          `json:"params"`
	Balances map[string]sdk.Int                    `json:"balances"`
	Deltas   map[string]sdk.Int                    `json:"deltas"`
	History  map[string][]types.RewardDistribution `json:"history"`
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for account := range data.Balances {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return err
		}
	}

	for account := range data.Deltas {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return err
		}
	}

	for account := range data.History {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return err
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:   types.DefaultParams(),
		Balances: make(map[string]sdk.Int),
		Deltas:   make(map[string]sdk.Int),
		History:  make(map[string][]types.RewardDistribution),
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for account, balance := range data.Balances {
		addr, _ := sdk.AccAddressFromBech32(account)
		keeper.SetBalance(ctx, addr, balance)
	}

	for account, delta := range data.Deltas {
		addr, _ := sdk.AccAddressFromBech32(account)
		keeper.IncBalanceDelta(ctx, addr, delta)
	}

	for account, history := range data.History {
		addr, _ := sdk.AccAddressFromBech32(account)
		for _, v := range history {
			keeper.AddRewardDistribution(ctx, addr, v.Height, v.Coins)
		}
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	balances := make(map[string]sdk.Int)
	deltas := make(map[string]sdk.Int)
	history := make(map[string][]types.RewardDistribution)

	it := keeper.GetBalanceIterator(ctx)
	for ; it.Valid(); it.Next() {
		owner := sdk.AccAddress(it.Key())
		balances[owner.String()] = keeper.GetBalance(ctx, it.Key())
	}
	it.Close()

	it = keeper.GetDeltasIterator(ctx)
	for ; it.Valid(); it.Next() {
		owner := sdk.AccAddress(it.Key())
		deltas[owner.String()] = keeper.GetBalanceDelta(ctx, it.Key())
	}
	it.Close()

	it = keeper.GetRewardsDistributionIterator(ctx)
	for ; it.Valid(); it.Next() {
		owner := sdk.AccAddress(it.Key())
		history[owner.String()] = keeper.GetRewardsDistributionHistory(ctx, it.Key())
	}
	it.Close()

	return GenesisState{
		Params:   keeper.GetParams(ctx),
		Balances: balances,
		Deltas:   deltas,
		History:  history,
	}
}
