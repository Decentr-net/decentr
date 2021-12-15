package v150

import (
	"fmt"

	communitytypes "github.com/Decentr-net/decentr/x/community/types"
	operationstypes "github.com/Decentr-net/decentr/x/operations/types"
	tokentypes "github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

type CommunityState struct {
	Params struct {
		Moderators []string `json:"moderators" yaml:"moderators"`
		FixedGas   struct {
			CreatePost sdk.Gas `json:"create_post" yaml:"create_post"`
			DeletePost sdk.Gas `json:"delete_post" yaml:"delete_post"`
			SetLike    sdk.Gas `json:"set_like" yaml:"set_like"`
			Follow     sdk.Gas `json:"follow" yaml:"follow"`
			Unfollow   sdk.Gas `json:"unfollow" yaml:"unfollow"`
		} `json:"fixed_gas" yaml:"fixed_gas"`
	} `json:"params"`
	Posts []struct {
		UUID         string         `json:"uuid"`
		Owner        sdk.AccAddress `json:"owner"`
		Title        string         `json:"title"`
		PreviewImage string         `json:"previewImage"`
		Category     uint8          `json:"category"`
		Text         string         `json:"text"`
	} `json:"posts"`
	Likes []struct {
		Owner     sdk.AccAddress `json:"owner"`
		PostOwner sdk.AccAddress `json:"postOwner"`
		PostUUID  string         `json:"postUuid"`
		Weight    int8           `json:"weight"`
	} `json:"likes"`
	Followers map[string][]string `json:"followers"`
}

type OperationsState struct {
	Params struct {
		Supervisors []string `json:"supervisors" yaml:"supervisors"`
		FixedGas    struct {
			ResetAccount      sdk.Gas `json:"reset_account" yaml:"reset_account"`
			BanAccount        sdk.Gas `json:"ban_account" yaml:"ban_account"`
			DistributeRewards sdk.Gas `json:"distribute_rewards" yaml:"distribute_rewards"`
		} `json:"fixed_gas" yaml:"fixed_gas"`
		MinGasPrice sdk.DecCoin `json:"min_gas_price" yaml:"min_gas_price"`
	} `json:"params"`
}

type TokenState struct {
	Params struct {
		RewardsBlockInterval int64 `json:"rewards_block_interval" yaml:"rewards_block_interval"`
	} `json:"params"`
	Balances map[string]sdk.Int `json:"balances"`
	Deltas   map[string]sdk.Int `json:"deltas"`
	History  map[string][]struct {
		Height int64     `json:"height"`
		Coins  sdk.Coins `json:"coins"`
	} `json:"history"`
}

func migrateDecentrAppState(appState genutiltypes.AppMap, clientCtx client.Context) genutiltypes.AppMap {
	v039codec, v040codec := clientCtx.LegacyAmino, clientCtx.Codec
	appState = migrateCommunity(appState, v039codec, v040codec)
	appState = migrateOperations(appState, v039codec, v040codec)
	appState = migrateToken(appState, v039codec, v040codec)

	return appState
}

func migrateCommunity(
	appState genutiltypes.AppMap,
	v039Codec *codec.LegacyAmino, v040Codec codec.Codec,
) genutiltypes.AppMap {
	if appState["community"] != nil {
		var oldState CommunityState
		v039Codec.MustUnmarshalJSON(appState["community"], &oldState)

		newState := communitytypes.DefaultGenesis()

		for _, v := range oldState.Params.Moderators {
			addr, err := sdk.AccAddressFromBech32(v)
			if err != nil {
				panic(fmt.Errorf("failed to parse community/moderators: %w", err))
			}
			newState.Params.Moderators = append(newState.Params.Moderators, addr)
		}
		newState.Params.FixedGas = communitytypes.FixedGasParams(oldState.Params.FixedGas)

		for _, v := range oldState.Posts {
			newState.Posts = append(newState.Posts, communitytypes.Post{
				Owner:        v.Owner,
				Uuid:         v.UUID,
				Title:        v.Title,
				PreviewImage: v.PreviewImage,
				Category:     communitytypes.Category(v.Category),
				Text:         v.Text,
			})
		}

		for _, v := range oldState.Likes {
			newState.Likes = append(newState.Likes, communitytypes.Like{
				Owner:     v.Owner,
				PostOwner: v.PostOwner,
				PostUuid:  v.PostUUID,
				Weight:    communitytypes.LikeWeight(v.Weight),
			})
		}

		for owner, followed := range oldState.Followers {
			var aa []sdk.AccAddress
			for _, v := range followed {
				addr, err := sdk.AccAddressFromBech32(v)
				if err != nil {
					panic(fmt.Errorf("failed to parse follower for %s: %w", owner, addr))
				}
				aa = append(aa, addr)
			}
			newState.Following[owner] = communitytypes.GenesisState_AddressList{Address: aa}
		}

		appState["community"] = v040Codec.MustMarshalJSON(newState)
	}

	return appState
}

func migrateOperations(
	appState genutiltypes.AppMap,
	v039Codec *codec.LegacyAmino, v040Codec codec.Codec,
) genutiltypes.AppMap {
	if appState["operations"] != nil {
		var oldState OperationsState
		v039Codec.MustUnmarshalJSON(appState["operations"], &oldState)

		newState := operationstypes.DefaultGenesis()

		for _, v := range oldState.Params.Supervisors {
			addr, err := sdk.AccAddressFromBech32(v)
			if err != nil {
				panic(fmt.Errorf("failed to parse operations/supervisors: %w", err))
			}
			newState.Params.Supervisors = append(newState.Params.Supervisors, addr)
		}
		newState.Params.FixedGas = operationstypes.FixedGasParams(oldState.Params.FixedGas)
		newState.Params.MinGasPrice = oldState.Params.MinGasPrice

		appState["operations"] = v040Codec.MustMarshalJSON(newState)
	}

	return appState
}

func migrateToken(
	appState genutiltypes.AppMap,
	v039Codec *codec.LegacyAmino, v040Codec codec.Codec,
) genutiltypes.AppMap {
	if appState["token"] != nil {
		var oldState TokenState
		v039Codec.MustUnmarshalJSON(appState["token"], &oldState)

		newState := tokentypes.DefaultGenesis()

		newState.Params.RewardsBlockInterval = uint64(oldState.Params.RewardsBlockInterval)
		for k, v := range oldState.Balances {
			newState.Balances[k] = sdk.DecProto{Dec: sdk.NewDecFromInt(v)}
		}
		for k, v := range oldState.Deltas {
			newState.Deltas[k] = sdk.DecProto{Dec: sdk.NewDecFromInt(v)}
		}

		appState["token"] = v040Codec.MustMarshalJSON(newState)
	}

	return appState
}
