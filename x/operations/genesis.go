package operations

import (
	"encoding/json"
	"fmt"

	"github.com/Decentr-net/decentr/x/operations/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Supervisors    []string       `json:"supervisors"`
	FixedGasParams FixedGasParams `json:"fixed_gas"`
	MinGasPrice    sdk.DecCoin    `json:"min_gas_price"`
}

// GetGenesisStateFromAppState returns community GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Supervisors:    types.DefaultSupervisors,
		FixedGasParams: types.DefaultFixedGasParams(),
		MinGasPrice:    types.DefaultMinGasPrice,
	}
}

// ValidateGenesis performs basic validation of PDV genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if len(data.Supervisors) == 0 {
		return fmt.Errorf("at least one supervisor has to be specified")
	}
	return nil
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetSupervisors(ctx, data.Supervisors)
	keeper.SetFixedGasParams(ctx, data.FixedGasParams)
	keeper.SetMinGasPrice(ctx, data.MinGasPrice)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Supervisors:    keeper.GetSupervisors(ctx),
		FixedGasParams: keeper.GetFixedGasParams(ctx),
		MinGasPrice:    keeper.GetMinGasPrice(ctx),
	}
}
