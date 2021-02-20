package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
)

var (
	DefaultCerberusOwners = []string{}
)

// ParamCerberusKey is store's key for ParamCerberus
var ParamCerberusOwnersKey = []byte("ParamCerberusOwners")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamCerberusOwnersKey, &DefaultCerberusOwners, validateCerberusOwners),
	)
}

func validateCerberusOwners(i interface{}) error {
	owners, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, owner := range owners {
		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			return fmt.Errorf("%s is an invalid cerberus address, err=%w", owner, err)
		}
	}
	return nil
}
