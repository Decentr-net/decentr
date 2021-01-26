package pdv

import (
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	CerberusAddr string `json:"cerberus_address"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(address string) GenesisState {
	return GenesisState{CerberusAddr: address}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState("https://cerberus.testnet.decentr.xyz")
}

// ValidateGenesis performs basic validation of PDV genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	_, err := url.ParseRequestURI(data.CerberusAddr)
	return err
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetCerberusAddr(ctx, data.CerberusAddr)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetCerberusAddr(ctx))
}
