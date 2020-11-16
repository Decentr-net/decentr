package community

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"

	"github.com/Decentr-net/decentr/x/community/types"
)

type GenesisState struct {
	PostRecords  []Post `json:"posts"`
	LikesRecords []Like `json:"likes"`
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.PostRecords {
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

	for _, record := range data.LikesRecords {
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

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		PostRecords:  []Post{},
		LikesRecords: []Like{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, record := range data.PostRecords {
		keeper.CreatePost(ctx, record)
	}

	for _, record := range data.LikesRecords {
		keeper.SetLike(ctx, record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var posts []Post
	iterator := k.GetPostsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		post := k.GetPostByKey(ctx, iterator.Key())
		posts = append(posts, post)
	}

	var likes []Like
	iterator = k.GetLikesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		like := k.GetLikeByKey(ctx, iterator.Key())
		likes = append(likes, like)
	}

	return GenesisState{PostRecords: posts, LikesRecords: likes}
}
