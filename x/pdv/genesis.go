package pdv

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	PDVRecords []PDV `json:"pdvs"`
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.PDVRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid PDVRecord: Address: %s. Error: Missing Owner", record.Address)
		}
		if record.Address == "" {
			return fmt.Errorf("invalid PDVRecord: Owner: %s. Error: Missing Address", record.Owner)
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
		keeper.SetPDV(ctx, record.Address, record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []PDV
	iterator := k.GetPDVsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		pdv := k.GetPDV(ctx, name)
		records = append(records, pdv)
	}
	return GenesisState{PDVRecords: records}
}
