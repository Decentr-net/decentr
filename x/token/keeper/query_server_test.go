package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/token/types"
)

func TestQueryServer_Balance(t *testing.T) {
	set, ctx := setupKeeper(t)
	s := NewQueryServer(set.keeper)

	address := NewAccAddress()
	set.keeper.IncTokens(ctx, address, sdk.OneDec())

	out, err := s.Balance(sdk.WrapSDKContext(ctx), &types.BalanceRequest{
		Address: address.String(),
	})
	require.NoError(t, err)
	require.Equal(t, &types.BalanceResponse{
		Balance: sdk.DecProto{Dec: sdk.NewDec(2)},
	}, out)
}
