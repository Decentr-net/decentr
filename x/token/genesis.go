package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Balances map[string]sdk.Int `json:"balances"`
}

func ValidateGenesis(data GenesisState) error {
	for account := range data.Balances {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return err
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Balances: make(map[string]sdk.Int),
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for account, balance := range data.Balances {
		addr, _ := sdk.AccAddressFromBech32(account)
		keeper.SetBalance(ctx, addr, balance)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	out := make(map[string]sdk.Int)

	it := keeper.GetBalanceIterator(ctx)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		amount := keeper.GetBalance(ctx, it.Key())
		owner := sdk.AccAddress(it.Key())
		out[owner.String()] = amount
	}

	return GenesisState{
		Balances: out,
	}
}
