package token

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token struct {
	Owner   sdk.AccAddress `json:"owner"`
	Balance sdk.Int        `json:"balance"`
}

type GenesisState struct {
	TokenRecords []Token `json:"tokens"`
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.TokenRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid Settings: Error: Missing Owner")
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		TokenRecords: []Token{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, record := range data.TokenRecords {
		// what we should do with index?
		keeper.AddTokens(ctx, record.Owner, time.Now(), record.Balance)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Token
	iterator := k.GetBalanceIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		owner := iterator.Key()
		balance := k.GetBalance(ctx, owner)
		records = append(records, Token{
			Owner:   owner,
			Balance: balance,
		})
	}
	return GenesisState{TokenRecords: records}
}
