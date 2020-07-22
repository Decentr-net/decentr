package decentr

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	PDVRecords []PDV `json:"pdv_records"`
}

func NewGenesisState(records []PDV) GenesisState {
	return GenesisState{PDVRecords: records}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.PDVRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid PDVRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid PDVRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		PDVRecords: []PDV{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, record := range data.PDVRecords {
		keeper.SetPDV(ctx, record.Value, record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []PDV
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		whois := k.GetPDV(ctx, name)
		records = append(records, whois)

	}
	return GenesisState{PDVRecords: records}
}
