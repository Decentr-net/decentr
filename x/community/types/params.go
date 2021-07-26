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
	DefaultModerators = make([]string, 0)
	DefaultFollowers  = make(map[string][]string)
)

var (
	KeyModerators = []byte("Moderators")
	KeyFixedGas   = []byte("FixedGas")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(KeyModerators, &DefaultModerators, validateModerators),
		params.NewParamSetPair(KeyFixedGas, FixedGasParams{}, validateFixedGasParams),
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
