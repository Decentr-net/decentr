package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{}
}
