package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/token/keeper"
	"github.com/Decentr-net/decentr/x/token/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {
	if keeper.SetParams(ctx, types.DefaultParams()); genState.Params != nil {
		keeper.SetParams(ctx, *genState.Params)
	}

	for k, v := range genState.Balances {
		address, err := sdk.AccAddressFromBech32(k)
		if err != nil {
			panic(fmt.Errorf("invalid address %s in balances : %w", k, err))
		}
		keeper.SetBalance(ctx, address, v.Dec)
	}

	accumulatedDelta := sdk.ZeroDec()
	for k, v := range genState.Deltas {
		address, err := sdk.AccAddressFromBech32(k)
		if err != nil {
			panic(fmt.Errorf("invalid address %s in deltas : %w", k, err))
		}
		keeper.SetBalanceDelta(ctx, address, v.Dec)
		accumulatedDelta = accumulatedDelta.Add(v.Dec)
	}
	keeper.IncAccumulatedDelta(ctx, accumulatedDelta)

	for _, v := range genState.BanList {
		address, err := sdk.AccAddressFromBech32(v)
		if err != nil {
			panic(fmt.Errorf("invalid address %s in ban list : %w", v, err))
		}
		keeper.SetBan(ctx, address, true)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	balances := map[string]sdk.DecProto{}
	deltas := map[string]sdk.DecProto{}
	banlist := make([]string, 0)

	k.IterateBalance(ctx, func(address sdk.AccAddress, balance sdk.Dec) (stop bool) {
		balances[address.String()] = sdk.DecProto{balance}
		return false
	})
	k.IterateBalanceDelta(ctx, func(address sdk.AccAddress, delta sdk.Dec) (stop bool) {
		deltas[address.String()] = sdk.DecProto{delta}
		return false
	})
	k.IterateBanList(ctx, func(address sdk.AccAddress) (stop bool) {
		banlist = append(banlist, address.String())
		return false
	})

	return &types.GenesisState{
		Params:   &params,
		Balances: balances,
		Deltas:   deltas,
		BanList:  banlist,
	}
}
