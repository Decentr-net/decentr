package operations_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Decentr-net/decentr/config"
	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/operations"
	"github.com/Decentr-net/decentr/x/operations/keeper"
	"github.com/Decentr-net/decentr/x/operations/types"
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
					Supervisors: types.DefaultParams().Supervisors,
					FixedGas:    types.DefaultParams().FixedGas,
					MinGasPrice: types.DefaultParams().MinGasPrice,
				},
			},
		},
		{
			name: "predefined",
			init: types.GenesisState{
				Params: &types.Params{
					Supervisors: []sdk.AccAddress{addr},
					FixedGas: types.FixedGasParams{
						ResetAccount:      1,
						DistributeRewards: 3,
					},
					MinGasPrice: sdk.NewDecCoin(config.DefaultBondDenom, sdk.NewInt(1)),
				},
			},
			exported: types.GenesisState{
				Params: &types.Params{
					Supervisors: []sdk.AccAddress{addr},
					FixedGas: types.FixedGasParams{
						ResetAccount:      1,
						DistributeRewards: 3,
					},
					MinGasPrice: sdk.NewDecCoin(config.DefaultBondDenom, sdk.NewInt(1)),
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

			operations.InitGenesis(ctx, *k, tc.init)
			got := operations.ExportGenesis(ctx, *k)
			require.Equal(t, tc.exported, *got)
		})
	}
}
