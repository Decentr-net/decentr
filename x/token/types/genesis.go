package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	p := DefaultParams()
	return &GenesisState{
		Params:   &p,
		Balances: map[string]sdk.DecProto{},
		Deltas:   map[string]sdk.DecProto{},
		BanList:  []sdk.AccAddress{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	validatePdvMap := func(m map[string]sdk.DecProto) error {
		for k, v := range m {
			if _, err := sdk.AccAddressFromBech32(k); err != nil {
				return fmt.Errorf("invalid address '%s': %w", k, err)
			}

			if v.Dec.IsNil() {
				return fmt.Errorf("invalid value for '%s': nil", k)
			}

			if v.Dec.IsNegative() {
				return fmt.Errorf("invalid value for '%s': negative value", k)
			}
		}
		return nil
	}

	if err := validatePdvMap(gs.Balances); err != nil {
		return fmt.Errorf("invalid balances: %w", err)
	}

	if err := validatePdvMap(gs.Deltas); err != nil {
		return fmt.Errorf("invalid deltas: %w", err)
	}

	for _, v := range gs.BanList {
		if err := sdk.VerifyAddressFormat(v); err != nil {
			return fmt.Errorf("invalid banned address '%s': %w", v, err)
		}
	}

	return nil
}
