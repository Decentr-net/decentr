package community

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"

	"github.com/Decentr-net/decentr/x/community/types"
)

type GenesisState struct {
	Posts          []Post              `json:"posts"`
	Likes          []Like              `json:"likes"`
	Moderators     []string            `json:"moderators"`
	FixedGasParams FixedGasParams      `json:"fixed_gas"`
	Followers      map[string][]string `json:"followers"`
}

// GetGenesisStateFromAppState returns community GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.Posts {
		if record.Owner == nil {
			return fmt.Errorf("invalid PostRecord: UUID: %s. Error: Missing Owner", record.UUID)
		}
		if record.UUID.Version() != uuid.V1 {
			return fmt.Errorf("invalid PostRecord: UUID: %s, Owner: %s. Error: Wrong UUID version", record.UUID, record.Owner)
		}
		if bytes.Equal(record.UUID.Bytes(), uuid.Nil.Bytes()) {
			return fmt.Errorf("invalid PostRecord: Owner: %s. Error: Empty UUID", record.Owner)
		}
		if record.Title == "" {
			return fmt.Errorf("invalid PostRecord: UUID: %s, Owner: %s. Error: Empty Title", record.UUID, record.Owner)
		}
		if record.Category == types.UndefinedCategory {
			return fmt.Errorf("invalid PostRecord: UUID: %s, Owner: %s. Error: Invalid Category", record.UUID, record.Owner)
		}
		if record.Text == "" {
			return fmt.Errorf("invalid PostRecord: UUID: %s, Owner: %s. Error: Empty Text", record.UUID, record.Owner)
		}
	}

	for _, record := range data.Likes {
		if record.Owner == nil {
			return fmt.Errorf("invalid LikeRecord: %+v. Error: Missing owner", record)
		}
		if record.PostOwner == nil {
			return fmt.Errorf("invalid LikeRecord: %+v. Error: Missing postOwner", record)
		}
		if record.PostUUID.Version() != uuid.V1 {
			return fmt.Errorf("invalid LikeRecord: %+v. Error: Wrong UUID version", record)
		}
		if record.Weight > types.LikeWeightUp || record.Weight < types.LikeWeightDown {
			return fmt.Errorf("invalid LikeRecord: %+v. Error: Invalid weight", record)
		}
	}

	if len(data.Moderators) == 0 {
		return fmt.Errorf("at least one moderator should be specified")
	}

	for who, whom := range data.Followers {
		if _, err := sdk.AccAddressFromBech32(who); err != nil {
			return err
		}
		for _, acc := range whom {
			if _, err := sdk.AccAddressFromBech32(acc); err != nil {
				return err
			}
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Posts:          []Post{},
		Likes:          []Like{},
		Moderators:     types.DefaultModerators,
		Followers:      types.DefaultFollowers,
		FixedGasParams: types.DefaultFixedGasParams(),
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, post := range data.Posts {
		keeper.CreatePost(ctx, post)
	}

	for _, like := range data.Likes {
		keeper.SetLike(ctx, like)
	}

	for who, whom := range data.Followers {
		whoAddr, _ := sdk.AccAddressFromBech32(who)
		for _, acc := range whom {
			whomAddr, _ := sdk.AccAddressFromBech32(acc)
			keeper.Follow(ctx, whoAddr, whomAddr)
		}
	}

	keeper.SetModerators(ctx, data.Moderators)
	keeper.SetFixedGasParams(ctx, data.FixedGasParams)
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var posts []Post

	it := k.GetPostsIterator(ctx)
	for ; it.Valid(); it.Next() {
		posts = append(posts, k.GetPostByKey(ctx, it.Key()))
	}
	it.Close()

	var likes []Like
	it = k.GetLikesIterator(ctx)
	for ; it.Valid(); it.Next() {
		likes = append(likes, k.GetLikeByKey(ctx, it.Key()))
	}

	var followers = make(map[string][]string)
	k.IterateFollowers(ctx, func(who, whom sdk.Address) (stop bool) {
		followers[who.String()] = append(followers[who.String()], whom.String())
		return false
	})

	return GenesisState{
		Posts:          posts,
		Likes:          likes,
		Followers:      followers,
		Moderators:     k.GetModerators(ctx),
		FixedGasParams: k.GetFixedGasParams(ctx),
	}
}
