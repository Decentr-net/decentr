package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultParamspace = ModuleName
)

var (
	DefaultModerators = []string(nil)
)

var (
	KeyModerators = []byte("Moderators")
	KeyFixedGas   = []byte("FixedGas")
)

// ParamKeyTable for operations module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyModerators, &p.Moderators, validateModerators),
		paramtypes.NewParamSetPair(KeyFixedGas, &p.FixedGas, validateFixedGasParams),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		Moderators: DefaultModerators,
		FixedGas:   DefaultFixedGasParams(),
	}
}

func validateModerators(i interface{}) error {
	owners, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, owner := range owners {
		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			return fmt.Errorf("%s is an invalid moderator address, err=%w", owner, err)
		}
	}
	return nil
}

func NewFixedGasParams(createPost, setLike, follow, unfollow sdk.Gas) FixedGasParams {
	return FixedGasParams{
		CretePost: createPost,
		SetLike:   setLike,
		Follow:    follow,
		Unfollow:  unfollow,
	}
}

func DefaultFixedGasParams() FixedGasParams {
	return NewFixedGasParams(0, 0, 0, 0)
}

func validateFixedGasParams(i interface{}) error {
	_, ok := i.(FixedGasParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateModerators(p.Moderators); err != nil {
		return fmt.Errorf("invalid moderators: %w", err)
	}

	if err := validateFixedGasParams(p.FixedGas); err != nil {
		return fmt.Errorf("invalid fixed_gas: %w", err)
	}

	return nil
}
