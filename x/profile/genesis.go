package profile

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	ProfileRecords []Profile `json:"settings_profile"`
}

func NewGenesisState(records []Profile) GenesisState {
	return GenesisState{ProfileRecords: records}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.ProfileRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid Settings: Error: Missing Owner")
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		ProfileRecords: []Profile{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, record := range data.ProfileRecords {
		keeper.SetProfile(ctx, record.Owner, record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Profile
	iterator := k.GetProfileIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := iterator.Key()
		whois := k.GetProfile(ctx, name)
		records = append(records, whois)

	}
	return GenesisState{ProfileRecords: records}
}