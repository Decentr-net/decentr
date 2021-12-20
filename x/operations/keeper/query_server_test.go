package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Decentr-net/decentr/x/operations/types"
)

func TestQueryServer_MinGasPrice(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewQueryServer(set.keeper, set.tokenKeeper)

	resp, err := s.MinGasPrice(sdk.WrapSDKContext(ctx), nil)
	require.NoError(t, err)
	require.Equal(t, types.DefaultMinGasPrice, resp.MinGasPrice)
}
