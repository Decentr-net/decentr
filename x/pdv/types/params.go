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

var (
	CerberusOwnersParamsKey = []byte("CerberusOwnersParams")
	FixedGasParamsKey       = []byte("FixedGasParams")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(CerberusOwnersParamsKey, &DefaultCerberusOwners, validateCerberusOwners),
		params.NewParamSetPair(FixedGasParamsKey, FixedGasParams{}, validateFixedGasParams),
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

type FixedGasParams struct {
	DeleteAccount     sdk.Gas `json:"delete_account" yaml:"delete_account"`
	DistributeRewards sdk.Gas `json:"distribute_rewards" yaml:"distribute_rewards"`
}

func NewFixedGasParams(deleteAccount, distributeReward sdk.Gas) FixedGasParams {
	return FixedGasParams{
		DeleteAccount:     deleteAccount,
		DistributeRewards: distributeReward,
	}
}

func DefaultFixedGasParams() FixedGasParams {
	return NewFixedGasParams(100, 100)
}

func validateFixedGasParams(i interface{}) error {
	v, ok := i.(FixedGasParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.DeleteAccount <= 0 {
		return fmt.Errorf("delete account be positive: %d", v.DeleteAccount)
	}

	if v.DistributeRewards <= 0 {
		return fmt.Errorf("distribute rewards be positive: %d", v.DistributeRewards)
	}

	return nil
}
