package types

import (
	"fmt"

	"github.com/Decentr-net/decentr/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultParamspace = ModuleName
)

var (
	DefaultSupervisors = []sdk.AccAddress(nil)
	DefaultMinGasPrice = sdk.NewDecCoinFromDec(config.DefaultBondDenom, sdk.MustNewDecFromStr("0.025"))
)

var (
	KeySupervisors = []byte("Supervisors")
	KeyFixedGas    = []byte("FixedGas")
	KeyMinGasPrice = []byte("MinGasPrice")
)

// ParamKeyTable for operations module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySupervisors, &p.Supervisors, validateSupervisors),
		paramtypes.NewParamSetPair(KeyFixedGas, &p.FixedGas, validateFixedGasParams),
		paramtypes.NewParamSetPair(KeyMinGasPrice, &p.MinGasPrice, validateMinGasPrice),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		Supervisors: DefaultSupervisors,
		FixedGas:    DefaultFixedGasParams(),
		MinGasPrice: DefaultMinGasPrice,
	}
}

func validateMinGasPrice(i interface{}) error {
	coin, ok := i.(sdk.DecCoin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !coin.IsValid() {
		return fmt.Errorf("coin is invalid")
	}

	if coin.IsNegative() {
		return fmt.Errorf("coin amount is negative")
	}

	if coin.IsZero() {
		return fmt.Errorf("coin amount is zero")
	}

	return nil
}

func validateSupervisors(i interface{}) error {
	s, ok := i.([]sdk.AccAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for i, v := range s {
		if err := sdk.VerifyAddressFormat(v); err != nil {
			return fmt.Errorf("invalid supervisor %d", i+1)
		}
	}

	return nil
}

func NewFixedGasParams(resetAccount, distributeReward, banAccount sdk.Gas) FixedGasParams {
	return FixedGasParams{
		ResetAccount:      resetAccount,
		DistributeRewards: distributeReward,
		BanAccount:        banAccount,
	}
}

func DefaultFixedGasParams() FixedGasParams {
	return NewFixedGasParams(0, 0, 0)
}

func validateFixedGasParams(i interface{}) error {
	_, ok := i.(FixedGasParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateSupervisors(p.Supervisors); err != nil {
		return fmt.Errorf("invalid supervisors: %w", err)
	}

	if err := validateFixedGasParams(p.FixedGas); err != nil {
		return fmt.Errorf("invalid fixed_gas: %w", err)
	}

	if err := validateMinGasPrice(p.MinGasPrice); err != nil {
		return fmt.Errorf("invalid min_gas_price: %w", err)
	}

	return nil
}
