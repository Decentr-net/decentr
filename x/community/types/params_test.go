package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
)

func Test_validateSupervisors(t *testing.T) {
	require.Error(t, validateModerators(nil))
	require.Error(t, validateModerators(1))
	require.NoError(t, validateModerators([]sdk.AccAddress{}))
	require.Error(t, validateModerators([]sdk.AccAddress{nil}))
	require.Error(t, validateModerators([]sdk.AccAddress{
		nil,
		NewAccAddress()},
	))
	require.NoError(t, validateModerators([]sdk.AccAddress{
		NewAccAddress()},
	))
	require.NoError(t, validateModerators([]sdk.AccAddress{
		NewAccAddress(),
		NewAccAddress(),
	}))
}

func Test_validateFixedGasParams(t *testing.T) {
	require.Error(t, validateFixedGasParams(nil))
	require.Error(t, validateFixedGasParams(1))
	require.NoError(t, validateFixedGasParams(FixedGasParams{}))
}
