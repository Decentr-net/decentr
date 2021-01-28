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
	DefaultModerators = []string{"decentr1nt5k6eg9zq5t2v66pr6zgyt5hh5tu8sk30re3a"}
)

// ParamModerators is store's key for ParamModerators
var ParamModeratorsKey = []byte("ParamModerators")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamModeratorsKey, &DefaultModerators, validateModerators),
	)
}

func validateModerators(i interface{}) error {
	moderators, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, moderator := range moderators {
		if _, err := sdk.AccAddressFromBech32(moderator); err != nil {
			return fmt.Errorf("%s is an invalid moderator address, err=%w", moderator, err)
		}
	}
	return nil
}
