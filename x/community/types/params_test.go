package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
)

func Test_validateModerators(t *testing.T) {
	require.Error(t, validateModerators(nil))
	require.Error(t, validateModerators(1))
	require.NoError(t, validateModerators([]string{}))
	require.Error(t, validateModerators([]string{""}))
	require.Error(t, validateModerators([]string{
		"",
		NewAccAddress().String(),
	}))
	require.NoError(t, validateModerators([]string{NewAccAddress().String()}))
	require.NoError(t, validateModerators([]string{
		NewAccAddress().String(),
		NewAccAddress().String(),
	}))
}

func Test_validateFixedGasParams(t *testing.T) {
	require.Error(t, validateFixedGasParams(nil))
	require.Error(t, validateFixedGasParams(1))
	require.NoError(t, validateFixedGasParams(FixedGasParams{}))
}
