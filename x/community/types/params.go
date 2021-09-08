package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
)

var (
	DefaultModerators = make([]string, 0)
	DefaultFollowers  = make(map[string][]string)
)

var (
	KeyModerators = []byte("Moderators")
	KeyFixedGas   = []byte("FixedGas")
)

type Params struct {
	Moderators []string       `json:"moderators" yaml:"moderators"`
	FixedGas   FixedGasParams `json:"fixed_gas" yaml:"fixed_gas"`
}

// ParamKeyTable for community module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(KeyModerators, &p.Moderators, validateModerators),
		params.NewParamSetPair(KeyFixedGas, &p.FixedGas, validateFixedGasParams),
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
	moderators, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(moderators) == 0 {
		return fmt.Errorf("can not be empty")
	}

	for _, moderator := range moderators {
		if _, err := sdk.AccAddressFromBech32(moderator); err != nil {
			return fmt.Errorf("%s is an invalid moderator address, err=%w", moderator, err)
		}
	}
	return nil
}

type FixedGasParams struct {
	CreatePost sdk.Gas `json:"create_post" yaml:"create_post"`
	DeletePost sdk.Gas `json:"delete_post" yaml:"delete_post"`
	SetLike    sdk.Gas `json:"set_like" yaml:"set_like"`
	Follow     sdk.Gas `json:"follow" yaml:"follow"`
	Unfollow   sdk.Gas `json:"unfollow" yaml:"unfollow"`
}

func DefaultFixedGasParams() FixedGasParams {
	return FixedGasParams{
		CreatePost: 1000,
		DeletePost: 100,
		SetLike:    100,
		Follow:     100,
		Unfollow:   100,
	}
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
