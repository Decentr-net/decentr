package token_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/token"
	"github.com/Decentr-net/decentr/x/token/keeper"
	"github.com/Decentr-net/decentr/x/token/types"
)

func TestGenesis(t *testing.T) {
	addr := NewAccAddress()

	tt := []struct {
		name     string
		init     types.GenesisState
		exported types.GenesisState
	}{
		{
			name: "default",
			exported: types.GenesisState{
				Balances: map[string]sdk.DecProto{},
			},
		},
		{
			name: "predefined",
			init: types.GenesisState{
				Balances: map[string]sdk.DecProto{
					addr.String(): {sdk.NewDec(1)},
				},
			},
			exported: types.GenesisState{
				Balances: map[string]sdk.DecProto{
					addr.String(): {sdk.NewDec(1)},
				},
			},
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			keys := sdk.NewKVStoreKeys(types.StoreKey, paramstypes.StoreKey)
			tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

			ctx, err := GetContext(keys, tkeys)
			require.NoError(t, err)

			registry := codectypes.NewInterfaceRegistry()
			cdc := codec.NewProtoCodec(registry)

			k := keeper.NewKeeper(
				cdc,
				keys[types.StoreKey],
			)

			token.InitGenesis(ctx, *k, tc.init)
			got := token.ExportGenesis(ctx, *k)
			require.Equal(t, tc.exported, *got)
		})
	}
}
