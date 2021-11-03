package token_test

import (
	"testing"

	"github.com/Decentr-net/decentr/x/token"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	. "github.com/Decentr-net/decentr/testutil"
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
				Params: &types.Params{
					RewardsBlockInterval: types.DefaultRewardsBlockInterval,
				},
				Balances: map[string]sdk.DecProto{},
				Deltas:   map[string]sdk.DecProto{},
				BanList:  []string{},
			},
		},
		{
			name: "predefined",
			init: types.GenesisState{
				Params: &types.Params{
					RewardsBlockInterval: types.DefaultRewardsBlockInterval,
				},
				Balances: map[string]sdk.DecProto{
					addr.String(): {sdk.NewDec(1)},
				},
				Deltas: map[string]sdk.DecProto{
					addr.String(): {sdk.NewDec(1)},
				},
				BanList: []string{addr.String()},
			},
			exported: types.GenesisState{
				Params: &types.Params{
					RewardsBlockInterval: types.DefaultRewardsBlockInterval,
				},
				Balances: map[string]sdk.DecProto{
					addr.String(): {sdk.NewDec(1)},
				},
				Deltas: map[string]sdk.DecProto{
					addr.String(): {sdk.NewDec(1)},
				},
				BanList: []string{addr.String()},
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

			paramsKeeper := paramskeeper.NewKeeper(
				cdc,
				codec.NewLegacyAmino(),
				keys[paramstypes.StoreKey],
				tkeys[paramstypes.TStoreKey],
			)

			k := keeper.NewKeeper(
				cdc,
				keys[types.StoreKey],
				paramsKeeper.Subspace(types.StoreKey),
			)

			token.InitGenesis(ctx, *k, tc.init)
			got := token.ExportGenesis(ctx, *k)
			require.Equal(t, tc.exported, *got)
		})
	}
}
