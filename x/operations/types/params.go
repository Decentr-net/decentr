package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName
	DefaultDenom      = "udec"
)

var (
	DefaultSupervisors = make([]string, 0)
	DefaultMinGasPrice = sdk.NewDecCoinFromDec(DefaultDenom, sdk.MustNewDecFromStr("0.025"))
)

var (
	KeySupervisors = []byte("Supervisors")
	KeyFixedGas    = []byte("FixedGas")
	KeyMinGasPrice = []byte("MinGasPrice")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(KeySupervisors, &DefaultSupervisors, validateSupervisors),
		params.NewParamSetPair(KeyFixedGas, FixedGasParams{}, validateFixedGasParams),
		params.NewParamSetPair(KeyMinGasPrice, &DefaultMinGasPrice, validateMinGasPrice),
	)
}

func validateMinGasPrice(i interface{}) error {
	coin, ok := i.(sdk.DecCoin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coin.Amount.IsZero() {
		return fmt.Errorf("amount cannot be zero")
	}

	return nil
}

func validateSupervisors(i interface{}) error {
	owners, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, owner := range owners {
		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			return fmt.Errorf("%s is an invalid supervisor address, err=%w", owner, err)
		}
	}
	return nil
}

type FixedGasParams struct {
	ResetAccount      sdk.Gas `json:"reset_account" yaml:"reset_account"`
	DistributeRewards sdk.Gas `json:"distribute_rewards" yaml:"distribute_rewards"`
}

func NewFixedGasParams(resetAccount, distributeReward sdk.Gas) FixedGasParams {
	return FixedGasParams{
		ResetAccount:      resetAccount,
		DistributeRewards: distributeReward,
	}
}

func DefaultFixedGasParams() FixedGasParams {
	return NewFixedGasParams(0, 0)
}

func validateFixedGasParams(i interface{}) error {
	_, ok := i.(FixedGasParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
