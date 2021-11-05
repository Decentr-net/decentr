package community_test

import (
	"testing"

	"github.com/Decentr-net/decentr/x/community"
	tokenkeeper "github.com/Decentr-net/decentr/x/token/keeper"
	tokentypes "github.com/Decentr-net/decentr/x/token/types"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	. "github.com/Decentr-net/decentr/testutil"
	"github.com/Decentr-net/decentr/x/community/keeper"
	"github.com/Decentr-net/decentr/x/community/types"
)

func TestGenesis(t *testing.T) {
	addr := NewAccAddress()
	postUUID := uuid.Must(uuid.NewV1())

	predefined := types.GenesisState{
		Params: &types.Params{
			Moderators: []string{addr.String()},
			FixedGas: types.FixedGasParams{
				CretePost: 1,
				SetLike:   2,
				Follow:    3,
				Unfollow:  4,
			},
		},
		Posts: []types.Post{
			{
				Uuid:         postUUID.String(),
				Owner:        addr.String(),
				Title:        "Title",
				PreviewImage: "",
				Category:     0,
				Text:         "Fifteen symbols should be typed",
			},
			{
				Uuid:         uuid.Must(uuid.NewV1()).String(),
				Owner:        addr.String(),
				Title:        "Title",
				PreviewImage: "",
				Category:     0,
				Text:         "Fifteen symbols should be typed",
			},
		},
		Likes: []types.Like{
			{
				Owner:     NewAccAddress().String(),
				PostOwner: addr.String(),
				PostUuid:  postUUID.String(),
				Weight:    1,
			},
		},
		Following: map[string]types.GenesisState_AddressList{
			NewAccAddress().String(): {
				[]string{NewAccAddress().String(), NewAccAddress().String()},
			},
			NewAccAddress().String(): {
				[]string{NewAccAddress().String(), NewAccAddress().String()},
			},
		},
	}

	tt := []struct {
		name     string
		init     types.GenesisState
		exported types.GenesisState
	}{
		{
			name: "default",
			exported: types.GenesisState{
				Params: &types.Params{
					Moderators: types.DefaultParams().Moderators,
					FixedGas:   types.DefaultParams().FixedGas,
				},
			},
		},
		{
			name:     "predefined",
			init:     predefined,
			exported: predefined,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			keys := sdk.NewKVStoreKeys(types.StoreKey, paramstypes.StoreKey, tokentypes.StoreKey)
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

			tokenKeeper := *tokenkeeper.NewKeeper(
				cdc,
				keys[tokentypes.StoreKey],
				paramsKeeper.Subspace(tokentypes.StoreKey),
			)
			tokenKeeper.SetParams(ctx, tokentypes.DefaultParams())

			k := keeper.NewKeeper(
				cdc,
				keys[types.StoreKey],
				paramsKeeper.Subspace(types.StoreKey),
			)

			community.InitGenesis(ctx, *k, tc.init)
			got := community.ExportGenesis(ctx, *k)
			require.NoError(t, got.Validate())

			require.Equal(t, tc.exported.Params, got.Params)
			require.ElementsMatch(t, tc.exported.Posts, got.Posts)
			require.ElementsMatch(t, tc.exported.Likes, got.Likes)

			require.Len(t, got.Following, len(tc.exported.Following))
			for k, v := range tc.exported.Following {
				require.ElementsMatch(t, v.Address, got.Following[k].Address)
			}
		})
	}
}
