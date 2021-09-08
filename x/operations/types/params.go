package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
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

type Params struct {
	Supervisors []string       `json:"supervisors" yaml:"supervisors"`
	FixedGas    FixedGasParams `json:"fixed_gas" yaml:"fixed_gas"`
	MinGasPrice sdk.DecCoin    `json:"min_gas_price" yaml:"min_gas_price"`
}

// ParamKeyTable for operations module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(KeySupervisors, &p.Supervisors, validateSupervisors),
		params.NewParamSetPair(KeyFixedGas, &p.FixedGas, validateFixedGasParams),
		params.NewParamSetPair(KeyMinGasPrice, &p.MinGasPrice, validateMinGasPrice),
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

	if len(owners) == 0 {
		return fmt.Errorf("can not be empty")
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
	BanAccount        sdk.Gas `json:"ban_account" yaml:"ban_account"`
	DistributeRewards sdk.Gas `json:"distribute_rewards" yaml:"distribute_rewards"`
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
