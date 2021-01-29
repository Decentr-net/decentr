package types

import (
	"fmt"
	"net/url"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
)

// ParamCerberusAddressKey is store's key for CerberusAddress
var ParamCerberusAddressKey = []byte("ParamCerberusAddress")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamCerberusAddressKey, "", validateCerberusAddress),
	)
}

func validateCerberusAddress(i interface{}) error {
	val, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := url.ParseRequestURI(val)
	return err
}
