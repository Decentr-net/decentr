package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

const (
	DefaultParamspace = ModuleName
)

var (
	DefaultRewardsBlockInterval int64 = 483840 // average block time is 5s, 17280 blocks a day, 28 days
)

var (
	KeyRewardsBlockInterval = []byte("RewardsBlockInterval")
)

type Params struct {
	RewardsBlockInterval int64 `json:"rewards_block_interval" yaml:"rewards_block_interval"`
}

// ParamKeyTable for token module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(KeyRewardsBlockInterval, &p.RewardsBlockInterval, validateRewardsBlockInterval),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		RewardsBlockInterval: DefaultRewardsBlockInterval,
	}
}

func validateRewardsBlockInterval(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return fmt.Errorf("invalid value: should be greater then 0")
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateRewardsBlockInterval(p.RewardsBlockInterval); err != nil {
		return fmt.Errorf("invalid rewards_block_interval: %w", err)
	}

	return nil
}
