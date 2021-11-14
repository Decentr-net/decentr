package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultParamspace = ModuleName
)

const (
	DefaultRewardsBlockInterval = 483840 // average block time is 5s, 17280 blocks a day, 28 days
)

var (
	KeyRewardsBlockInterval = []byte("RewardsBlockInterval")
)

// ParamKeyTable for operations module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRewardsBlockInterval, &p.RewardsBlockInterval, validateRewardsBlockInterval),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		RewardsBlockInterval: DefaultRewardsBlockInterval,
	}
}

func validateRewardsBlockInterval(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("interval is zero")
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateRewardsBlockInterval(p.RewardsBlockInterval); err != nil {
		return fmt.Errorf("invalid rewards_block_interval: %w", err)
	}

	return nil
}
