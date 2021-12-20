package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
)

func Test_validateSupervisors(t *testing.T) {
	require.Error(t, validateSupervisors(nil))
	require.Error(t, validateSupervisors(1))
	require.NoError(t, validateSupervisors([]string{}))
	require.Error(t, validateSupervisors([]string{
		"",
		NewAccAddress().String(),
	}))
	require.NoError(t, validateSupervisors([]string{
		NewAccAddress().String(),
	}))
	require.NoError(t, validateSupervisors([]string{
		NewAccAddress().String(),
		NewAccAddress().String(),
	}))
}

func Test_validateMinGasPrice(t *testing.T) {
	require.Error(t, validateMinGasPrice(nil))
	require.Error(t, validateMinGasPrice(1))
	require.Error(t, validateMinGasPrice(sdk.DecCoin{}))
	require.Error(t, validateMinGasPrice(sdk.DecCoin{
		Denom:  "",
		Amount: sdk.NewDec(1),
	}))
	require.Error(t, validateMinGasPrice(sdk.DecCoin{
		Denom:  config.DefaultBondDenom,
		Amount: sdk.NewDec(-1),
	}))
	require.Error(t, validateMinGasPrice(sdk.DecCoin{
		Denom:  config.DefaultBondDenom,
		Amount: sdk.NewDec(0),
	}))
	require.NoError(t, validateMinGasPrice(sdk.DecCoin{
		Denom:  config.DefaultBondDenom,
		Amount: sdk.NewDec(1),
	}))
}

func Test_validateFixedGasParams(t *testing.T) {
	require.Error(t, validateFixedGasParams(nil))
	require.Error(t, validateFixedGasParams(1))
	require.NoError(t, validateFixedGasParams(FixedGasParams{}))
}
