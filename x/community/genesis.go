package community

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"
)

type GenesisState struct {
	PostRecords []Post
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
		if record.Text == "" {
			return fmt.Errorf("invalid PostRecord: UUID: %s, Owner: %s. Error: Empty Text", record.UUID, record.Owner)
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		PostRecords: []Post{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, record := range data.PostRecords {
		keeper.CreatePost(ctx, record.Owner, record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Post
	iterator := k.GetPostsIterator(ctx, "")
	for ; iterator.Valid(); iterator.Next() {
		post := k.GetPostByKey(ctx, iterator.Key())
		records = append(records, post)
	}

	return GenesisState{PostRecords: records}
}
